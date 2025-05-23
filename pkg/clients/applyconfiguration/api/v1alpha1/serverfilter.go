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

// ServerFilterApplyConfiguration represents a declarative configuration of the ServerFilter type for use
// with apply.
type ServerFilterApplyConfiguration struct {
	Name                                 *apiv1alpha1.OpenStackName `json:"name,omitempty"`
	FilterByServerTagsApplyConfiguration `json:",inline"`
}

// ServerFilterApplyConfiguration constructs a declarative configuration of the ServerFilter type for use with
// apply.
func ServerFilter() *ServerFilterApplyConfiguration {
	return &ServerFilterApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ServerFilterApplyConfiguration) WithName(value apiv1alpha1.OpenStackName) *ServerFilterApplyConfiguration {
	b.Name = &value
	return b
}

// WithTags adds the given value to the Tags field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tags field.
func (b *ServerFilterApplyConfiguration) WithTags(values ...apiv1alpha1.ServerTag) *ServerFilterApplyConfiguration {
	for i := range values {
		b.FilterByServerTagsApplyConfiguration.Tags = append(b.FilterByServerTagsApplyConfiguration.Tags, values[i])
	}
	return b
}

// WithTagsAny adds the given value to the TagsAny field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the TagsAny field.
func (b *ServerFilterApplyConfiguration) WithTagsAny(values ...apiv1alpha1.ServerTag) *ServerFilterApplyConfiguration {
	for i := range values {
		b.FilterByServerTagsApplyConfiguration.TagsAny = append(b.FilterByServerTagsApplyConfiguration.TagsAny, values[i])
	}
	return b
}

// WithNotTags adds the given value to the NotTags field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NotTags field.
func (b *ServerFilterApplyConfiguration) WithNotTags(values ...apiv1alpha1.ServerTag) *ServerFilterApplyConfiguration {
	for i := range values {
		b.FilterByServerTagsApplyConfiguration.NotTags = append(b.FilterByServerTagsApplyConfiguration.NotTags, values[i])
	}
	return b
}

// WithNotTagsAny adds the given value to the NotTagsAny field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NotTagsAny field.
func (b *ServerFilterApplyConfiguration) WithNotTagsAny(values ...apiv1alpha1.ServerTag) *ServerFilterApplyConfiguration {
	for i := range values {
		b.FilterByServerTagsApplyConfiguration.NotTagsAny = append(b.FilterByServerTagsApplyConfiguration.NotTagsAny, values[i])
	}
	return b
}
