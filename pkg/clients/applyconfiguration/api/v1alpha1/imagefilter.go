/*
Copyright 2024 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
)

// ImageFilterApplyConfiguration represents a declarative configuration of the ImageFilter type for use
// with apply.
type ImageFilterApplyConfiguration struct {
	Name       *apiv1alpha1.OpenStackName   `json:"name,omitempty"`
	Visibility *apiv1alpha1.ImageVisibility `json:"visibility,omitempty"`
	Tags       []apiv1alpha1.ImageTag       `json:"tags,omitempty"`
}

// ImageFilterApplyConfiguration constructs a declarative configuration of the ImageFilter type for use with
// apply.
func ImageFilter() *ImageFilterApplyConfiguration {
	return &ImageFilterApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ImageFilterApplyConfiguration) WithName(value apiv1alpha1.OpenStackName) *ImageFilterApplyConfiguration {
	b.Name = &value
	return b
}

// WithVisibility sets the Visibility field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Visibility field is set to the value of the last call.
func (b *ImageFilterApplyConfiguration) WithVisibility(value apiv1alpha1.ImageVisibility) *ImageFilterApplyConfiguration {
	b.Visibility = &value
	return b
}

// WithTags adds the given value to the Tags field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tags field.
func (b *ImageFilterApplyConfiguration) WithTags(values ...apiv1alpha1.ImageTag) *ImageFilterApplyConfiguration {
	for i := range values {
		b.Tags = append(b.Tags, values[i])
	}
	return b
}
