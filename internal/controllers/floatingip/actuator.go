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

package floatingip

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/neutrontags"
)

type (
	osResourceT = floatingips.FloatingIP

	createResourceActuator    = interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT]
	deleteResourceActuator    = interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT]
	reconcileResourceActuator = interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]
	resourceReconciler        = interfaces.ResourceReconciler[orcObjectPT, osResourceT]
	helperFactory             = interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	floatingipIterator        = iter.Seq2[*osResourceT, error]
)

type floatingipActuator struct {
	osClient osclients.NetworkClient
}

type floatingipCreateActuator struct {
	floatingipActuator
	k8sClient client.Client
}

var _ createResourceActuator = floatingipCreateActuator{}
var _ deleteResourceActuator = floatingipActuator{}

func (floatingipActuator) GetResourceID(osResource *floatingips.FloatingIP) string {
	return osResource.ID
}

func (actuator floatingipActuator) GetOSResourceByID(ctx context.Context, id string) (*floatingips.FloatingIP, progress.ReconcileStatus) {
	floatingip, err := actuator.osClient.GetFloatingIP(ctx, id)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	return floatingip, nil
}

func (actuator floatingipActuator) ListOSResourcesForAdoption(ctx context.Context, obj *orcv1alpha1.FloatingIP) (floatingipIterator, bool) {
	if obj.Spec.Resource == nil {
		return nil, false
	}

	listOpts := floatingips.ListOpts{}
	return actuator.osClient.ListFloatingIP(ctx, listOpts), true
}

func (actuator floatingipCreateActuator) ListOSResourcesForImport(ctx context.Context, obj orcObjectPT, filter filterT) (iter.Seq2[*osResourceT, error], progress.ReconcileStatus) {
	listOpts := floatingips.ListOpts{
		Description: string(ptr.Deref(filter.Description, "")),
		Tags:        neutrontags.Join(filter.Tags),
		TagsAny:     neutrontags.Join(filter.TagsAny),
		NotTags:     neutrontags.Join(filter.NotTags),
		NotTagsAny:  neutrontags.Join(filter.NotTagsAny),
	}

	return actuator.osClient.ListFloatingIP(ctx, listOpts), nil
}

func (actuator floatingipCreateActuator) CreateResource(ctx context.Context, obj *orcv1alpha1.FloatingIP) (*floatingips.FloatingIP, progress.ReconcileStatus) {
	resource := obj.Spec.Resource
	if resource == nil {
		// Should have been caught by API validation
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Creation requested, but spec.resource is not set"))
	}

	var floatingNetworkID *string
	{
		// Fetch dependencies and ensure they have our finalizer
		externalNetwork, reconcileStatus := externalNetworkDep.GetDependency(
			ctx, actuator.k8sClient, obj, func(dep *orcv1alpha1.Network) bool {
				return orcv1alpha1.IsAvailable(dep) && dep.Status.ID != nil
			},
		)
		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
			return nil, reconcileStatus
		}
		floatingNetworkID = externalNetwork.Status.ID
	}

	createOpts := floatingips.CreateOpts{
		Description:       string(ptr.Deref(resource.Description, "")),
		FloatingNetworkID: *floatingNetworkID,
	}
	if resource.FloatingIP != nil {
		createOpts.FloatingIP = string(*resource.FloatingIP)
	}

	osResource, err := actuator.osClient.CreateFloatingIP(ctx, &createOpts)

	// We should require the spec to be updated before retrying a create which returned a conflict
	if orcerrors.IsConflict(err) {
		err = orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "invalid configuration creating resource: "+err.Error(), err)
	}

	if err != nil {
		return nil, progress.WrapError(err)
	}
	return osResource, nil
}

func (actuator floatingipActuator) DeleteResource(ctx context.Context, _ orcObjectPT, floatingip *floatingips.FloatingIP) progress.ReconcileStatus {
	return progress.WrapError(actuator.osClient.DeleteFloatingIP(ctx, floatingip.ID))
}

var _ reconcileResourceActuator = floatingipActuator{}

func (actuator floatingipActuator) GetResourceReconcilers(ctx context.Context, orcObject orcObjectPT, osResource *osResourceT, controller interfaces.ResourceController) ([]resourceReconciler, progress.ReconcileStatus) {
	return []resourceReconciler{
		neutrontags.ReconcileTags[orcObjectPT, osResourceT](actuator.osClient, "floatingips", osResource.ID, orcObject.Spec.Resource.Tags, osResource.Tags),
	}, nil
}

type floatingipHelperFactory struct{}

var _ helperFactory = floatingipHelperFactory{}

func (floatingipHelperFactory) NewAPIObjectAdapter(obj orcObjectPT) adapterI {
	return floatingipAdapter{obj}
}

func (floatingipHelperFactory) NewCreateActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (createResourceActuator, progress.ReconcileStatus) {
	return newCreateActuator(ctx, orcObject, controller)
}

func (floatingipHelperFactory) NewDeleteActuator(ctx context.Context, orcObject orcObjectPT, controller interfaces.ResourceController) (deleteResourceActuator, progress.ReconcileStatus) {
	return newActuator(ctx, orcObject, controller)
}

func newActuator(ctx context.Context, orcObject *orcv1alpha1.FloatingIP, controller interfaces.ResourceController) (floatingipActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, controller.GetK8sClient(), orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return floatingipActuator{}, reconcileStatus
	}

	clientScope, err := controller.GetScopeFactory().NewClientScopeFromObject(ctx, controller.GetK8sClient(), log, orcObject)
	if err != nil {
		return floatingipActuator{}, nil
	}
	osClient, err := clientScope.NewNetworkClient()
	if err != nil {
		return floatingipActuator{}, nil
	}

	return floatingipActuator{
		osClient: osClient,
	}, nil
}

func newCreateActuator(ctx context.Context, orcObject *orcv1alpha1.FloatingIP, controller interfaces.ResourceController) (floatingipCreateActuator, progress.ReconcileStatus) {
	floatingipActuator, reconcileStatus := newActuator(ctx, orcObject, controller)
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return floatingipCreateActuator{}, reconcileStatus
	}

	return floatingipCreateActuator{
		floatingipActuator: floatingipActuator,
		k8sClient:          controller.GetK8sClient(),
	}, nil
}
