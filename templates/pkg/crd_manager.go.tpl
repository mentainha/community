{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"
	"fmt"
	ackerr "github.com/aws/aws-controllers-k8s/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	ackv1alpha1 "github.com/aws/aws-controllers-k8s/apis/core/v1alpha1"
	ackcompare "github.com/aws/aws-controllers-k8s/pkg/compare"
	ackmetrics "github.com/aws/aws-controllers-k8s/pkg/metrics"
	acktypes "github.com/aws/aws-controllers-k8s/pkg/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"

	svcsdk "github.com/aws/aws-sdk-go/service/{{ .ServiceIDClean }}"
	svcsdkapi "github.com/aws/aws-sdk-go/service/{{ .ServiceIDClean }}/{{ .ServiceIDClean }}iface"
)

// +kubebuilder:rbac:groups={{ .APIGroup }},resources={{ ToLower .CRD.Plural }},verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{ .APIGroup }},resources={{ ToLower .CRD.Plural }}/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// resourceManager is responsible for providing a consistent way to perform
// CRUD operations in a backend AWS service API for Book custom resources.
type resourceManager struct {
	// log refers to the logr.Logger object handling logging for the service
	// controller
	log logr.Logger
	// metrics contains a collection of Prometheus metric objects that the
	// service controller and its reconcilers track
	metrics *ackmetrics.Metrics
	// rr is the AWSResourceReconciler which can be used for various utility
	// functions such as querying for Secret values given a SecretReference
	rr acktypes.AWSResourceReconciler
	// awsAccountID is the AWS account identifier that contains the resources
	// managed by this resource manager
	awsAccountID ackv1alpha1.AWSAccountID
	// The AWS Region that this resource manager targets
	awsRegion ackv1alpha1.AWSRegion
	// sess is the AWS SDK Session object used to communicate with the backend
	// AWS service API
	sess *session.Session
	// sdk is a pointer to the AWS service API interface exposed by the
	// aws-sdk-go/services/{alias}/{alias}iface package.
	sdkapi svcsdkapi.{{ .SDKAPIInterfaceTypeName }}API
}

// concreteResource returns a pointer to a resource from the supplied
// generic AWSResource interface
func (rm *resourceManager) concreteResource(
	res acktypes.AWSResource,
) *resource {
	// cast the generic interface into a pointer type specific to the concrete
	// implementing resource type managed by this resource manager
	return res.(*resource)
}

// ReadOne returns the currently-observed state of the supplied AWSResource in
// the backend AWS service API.
func (rm *resourceManager) ReadOne(
	ctx context.Context,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's ReadOne() method received resource with nil CR object")
	}
	observed, err := rm.sdkFind(ctx, r)
	if err != nil {
		return nil, err
	}
	return rm.onSuccess(observed)
}

// Create attempts to create the supplied AWSResource in the backend AWS
// service API, returning an AWSResource representing the newly-created
// resource
func (rm *resourceManager) Create(
	ctx context.Context,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Create() method received resource with nil CR object")
	}
	created, err := rm.sdkCreate(ctx, r)
	if err != nil {
		return rm.onError(r, err)
	}
	return rm.onSuccess(created)
}

// Update attempts to mutate the supplied desired AWSResource in the backend AWS
// service API, returning an AWSResource representing the newly-mutated
// resource.
// Note for specialized logic implementers can check to see how the latest
// observed resource differs from the supplied desired state. The
// higher-level reonciler determines whether or not the desired differs
// from the latest observed and decides whether to call the resource
// manager's Update method
func (rm *resourceManager) Update(
	ctx context.Context,
	resDesired acktypes.AWSResource,
	resLatest acktypes.AWSResource,
	diffReporter *ackcompare.Reporter,
) (acktypes.AWSResource, error) {
	desired := rm.concreteResource(resDesired)
	latest := rm.concreteResource(resLatest)
	if desired.ko == nil || latest.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Update() method received resource with nil CR object")
	}
	updated, err := rm.sdkUpdate(ctx, desired, latest, diffReporter)
	if err != nil {
		return rm.onError(latest, err)
	}
	return rm.onSuccess(updated)
}

// Delete attempts to destroy the supplied AWSResource in the backend AWS
// service API.
func (rm *resourceManager) Delete(
	ctx context.Context,
	res acktypes.AWSResource,
) error {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Update() method received resource with nil CR object")
	}
	return rm.sdkDelete(ctx, r)
}

// ARNFromName returns an AWS Resource Name from a given string name. This
// is useful for constructing ARNs for APIs that require ARNs in their
// GetAttributes operations but all we have (for new CRs at least) is a
// name for the resource
func (rm *resourceManager) ARNFromName(name string) string {
	return fmt.Sprintf(
		"arn:aws:{{ .ServiceIDClean }}:%s:%s:%s",
		rm.awsRegion,
		rm.awsAccountID,
		name,
	)
}

// newResourceManager returns a new struct implementing
// acktypes.AWSResourceManager
func newResourceManager(
	log logr.Logger,
	metrics *ackmetrics.Metrics,
	rr acktypes.AWSResourceReconciler,
	sess *session.Session,
	id ackv1alpha1.AWSAccountID,
	region ackv1alpha1.AWSRegion,
) (*resourceManager, error) {
	return &resourceManager{
		log: log,
		metrics: metrics,
		rr: rr,
		awsAccountID: id,
		awsRegion: region,
		sess:		 sess,
		sdkapi:	   svcsdk.New(sess),
	}, nil
}

// onError updates resource conditions and returns updated resource
// it returns nil if no condition is updated.
func (rm *resourceManager) onError(
	r *resource,
	err error,
) (*resource, error) {
	r1, updated := rm.updateConditions(r, err)
	if !updated {
		return nil, err
	}
	for _, condition := range r1.Conditions() {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal &&
			condition.Status == corev1.ConditionTrue {
			// resource is in Terminal condition
			// return Terminal error
			return r1, ackerr.Terminal
		}
	}
	return r1, err
}

// onSuccess updates resource conditions and returns updated resource
// it returns the supplied resource if no condition is updated.
func (rm *resourceManager) onSuccess(
	r *resource,
) (*resource, error) {
	r1, updated := rm.updateConditions(r, nil)
	if !updated {
		return r, nil
	}
	return r1, nil
}
