// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package authorizer

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	ackv1alpha1 "github.com/aws/aws-controllers-k8s/apis/core/v1alpha1"
	ackcompare "github.com/aws/aws-controllers-k8s/pkg/compare"
	ackerr "github.com/aws/aws-controllers-k8s/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/apigatewayv2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws/aws-controllers-k8s/services/apigatewayv2/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = &aws.JSONValue{}
	_ = &svcsdk.ApiGatewayV2{}
	_ = &svcapitypes.Authorizer{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.GetAuthorizerWithContext(ctx, input)
	if respErr != nil {
		if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "NotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.AuthorizerCredentialsArn != nil {
		ko.Spec.AuthorizerCredentialsARN = resp.AuthorizerCredentialsArn
	}
	if resp.AuthorizerId != nil {
		ko.Status.AuthorizerID = resp.AuthorizerId
	}
	if resp.AuthorizerPayloadFormatVersion != nil {
		ko.Spec.AuthorizerPayloadFormatVersion = resp.AuthorizerPayloadFormatVersion
	}
	if resp.AuthorizerResultTtlInSeconds != nil {
		ko.Spec.AuthorizerResultTtlInSeconds = resp.AuthorizerResultTtlInSeconds
	}
	if resp.AuthorizerType != nil {
		ko.Spec.AuthorizerType = resp.AuthorizerType
	}
	if resp.AuthorizerUri != nil {
		ko.Spec.AuthorizerURI = resp.AuthorizerUri
	}
	if resp.EnableSimpleResponses != nil {
		ko.Spec.EnableSimpleResponses = resp.EnableSimpleResponses
	}
	if resp.IdentitySource != nil {
		f7 := []*string{}
		for _, f7iter := range resp.IdentitySource {
			var f7elem string
			f7elem = *f7iter
			f7 = append(f7, &f7elem)
		}
		ko.Spec.IDentitySource = f7
	}
	if resp.IdentityValidationExpression != nil {
		ko.Spec.IDentityValidationExpression = resp.IdentityValidationExpression
	}
	if resp.JwtConfiguration != nil {
		f9 := &svcapitypes.JWTConfiguration{}
		if resp.JwtConfiguration.Audience != nil {
			f9f0 := []*string{}
			for _, f9f0iter := range resp.JwtConfiguration.Audience {
				var f9f0elem string
				f9f0elem = *f9f0iter
				f9f0 = append(f9f0, &f9f0elem)
			}
			f9.Audience = f9f0
		}
		if resp.JwtConfiguration.Issuer != nil {
			f9.Issuer = resp.JwtConfiguration.Issuer
		}
		ko.Spec.JWTConfiguration = f9
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required by not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Status.AuthorizerID == nil || r.ko.Spec.APIID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetAuthorizerInput, error) {
	res := &svcsdk.GetAuthorizerInput{}

	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}
	if r.ko.Status.AuthorizerID != nil {
		res.SetAuthorizerId(*r.ko.Status.AuthorizerID)
	}

	return res, nil
}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.GetAuthorizersInput, error) {
	res := &svcsdk.GetAuthorizersInput{}

	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a new resource with any fields in the Status field filled in
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	input, err := rm.newCreateRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.CreateAuthorizerWithContext(ctx, input)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.AuthorizerId != nil {
		ko.Status.AuthorizerID = resp.AuthorizerId
	}

	rm.setStatusDefaults(ko)

	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	r *resource,
) (*svcsdk.CreateAuthorizerInput, error) {
	res := &svcsdk.CreateAuthorizerInput{}

	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}
	if r.ko.Spec.AuthorizerCredentialsARN != nil {
		res.SetAuthorizerCredentialsArn(*r.ko.Spec.AuthorizerCredentialsARN)
	}
	if r.ko.Spec.AuthorizerPayloadFormatVersion != nil {
		res.SetAuthorizerPayloadFormatVersion(*r.ko.Spec.AuthorizerPayloadFormatVersion)
	}
	if r.ko.Spec.AuthorizerResultTtlInSeconds != nil {
		res.SetAuthorizerResultTtlInSeconds(*r.ko.Spec.AuthorizerResultTtlInSeconds)
	}
	if r.ko.Spec.AuthorizerType != nil {
		res.SetAuthorizerType(*r.ko.Spec.AuthorizerType)
	}
	if r.ko.Spec.AuthorizerURI != nil {
		res.SetAuthorizerUri(*r.ko.Spec.AuthorizerURI)
	}
	if r.ko.Spec.EnableSimpleResponses != nil {
		res.SetEnableSimpleResponses(*r.ko.Spec.EnableSimpleResponses)
	}
	if r.ko.Spec.IDentitySource != nil {
		f7 := []*string{}
		for _, f7iter := range r.ko.Spec.IDentitySource {
			var f7elem string
			f7elem = *f7iter
			f7 = append(f7, &f7elem)
		}
		res.SetIdentitySource(f7)
	}
	if r.ko.Spec.IDentityValidationExpression != nil {
		res.SetIdentityValidationExpression(*r.ko.Spec.IDentityValidationExpression)
	}
	if r.ko.Spec.JWTConfiguration != nil {
		f9 := &svcsdk.JWTConfiguration{}
		if r.ko.Spec.JWTConfiguration.Audience != nil {
			f9f0 := []*string{}
			for _, f9f0iter := range r.ko.Spec.JWTConfiguration.Audience {
				var f9f0elem string
				f9f0elem = *f9f0iter
				f9f0 = append(f9f0, &f9f0elem)
			}
			f9.SetAudience(f9f0)
		}
		if r.ko.Spec.JWTConfiguration.Issuer != nil {
			f9.SetIssuer(*r.ko.Spec.JWTConfiguration.Issuer)
		}
		res.SetJwtConfiguration(f9)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	diffReporter *ackcompare.Reporter,
) (*resource, error) {

	input, err := rm.newUpdateRequestPayload(desired)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.UpdateAuthorizerWithContext(ctx, input)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.AuthorizerId != nil {
		ko.Status.AuthorizerID = resp.AuthorizerId
	}

	rm.setStatusDefaults(ko)

	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	r *resource,
) (*svcsdk.UpdateAuthorizerInput, error) {
	res := &svcsdk.UpdateAuthorizerInput{}

	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}
	if r.ko.Spec.AuthorizerCredentialsARN != nil {
		res.SetAuthorizerCredentialsArn(*r.ko.Spec.AuthorizerCredentialsARN)
	}
	if r.ko.Status.AuthorizerID != nil {
		res.SetAuthorizerId(*r.ko.Status.AuthorizerID)
	}
	if r.ko.Spec.AuthorizerPayloadFormatVersion != nil {
		res.SetAuthorizerPayloadFormatVersion(*r.ko.Spec.AuthorizerPayloadFormatVersion)
	}
	if r.ko.Spec.AuthorizerResultTtlInSeconds != nil {
		res.SetAuthorizerResultTtlInSeconds(*r.ko.Spec.AuthorizerResultTtlInSeconds)
	}
	if r.ko.Spec.AuthorizerType != nil {
		res.SetAuthorizerType(*r.ko.Spec.AuthorizerType)
	}
	if r.ko.Spec.AuthorizerURI != nil {
		res.SetAuthorizerUri(*r.ko.Spec.AuthorizerURI)
	}
	if r.ko.Spec.EnableSimpleResponses != nil {
		res.SetEnableSimpleResponses(*r.ko.Spec.EnableSimpleResponses)
	}
	if r.ko.Spec.IDentitySource != nil {
		f8 := []*string{}
		for _, f8iter := range r.ko.Spec.IDentitySource {
			var f8elem string
			f8elem = *f8iter
			f8 = append(f8, &f8elem)
		}
		res.SetIdentitySource(f8)
	}
	if r.ko.Spec.IDentityValidationExpression != nil {
		res.SetIdentityValidationExpression(*r.ko.Spec.IDentityValidationExpression)
	}
	if r.ko.Spec.JWTConfiguration != nil {
		f10 := &svcsdk.JWTConfiguration{}
		if r.ko.Spec.JWTConfiguration.Audience != nil {
			f10f0 := []*string{}
			for _, f10f0iter := range r.ko.Spec.JWTConfiguration.Audience {
				var f10f0elem string
				f10f0elem = *f10f0iter
				f10f0 = append(f10f0, &f10f0elem)
			}
			f10.SetAudience(f10f0)
		}
		if r.ko.Spec.JWTConfiguration.Issuer != nil {
			f10.SetIssuer(*r.ko.Spec.JWTConfiguration.Issuer)
		}
		res.SetJwtConfiguration(f10)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) error {
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return err
	}
	_, respErr := rm.sdkapi.DeleteAuthorizerWithContext(ctx, input)
	return respErr
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteAuthorizerInput, error) {
	res := &svcsdk.DeleteAuthorizerInput{}

	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}
	if r.ko.Status.AuthorizerID != nil {
		res.SetAuthorizerId(*r.ko.Status.AuthorizerID)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Authorizer,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
			break
		}
	}

	if rm.terminalAWSError(err) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		terminalCondition.Status = corev1.ConditionTrue
		awsErr, _ := ackerr.AWSError(err)
		errorMessage := awsErr.Message()
		terminalCondition.Message = &errorMessage
	} else if terminalCondition != nil {
		terminalCondition.Status = corev1.ConditionFalse
		terminalCondition.Message = nil
	}
	if terminalCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
