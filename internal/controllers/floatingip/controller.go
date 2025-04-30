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
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/pkg/predicates"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/credentials"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/dependency"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=floatingips,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openstack.k-orc.cloud,resources=floatingips/status,verbs=get;update;patch

type floatingipReconcilerConstructor struct {
	scopeFactory scope.Factory
}

func New(scopeFactory scope.Factory) interfaces.Controller {
	return floatingipReconcilerConstructor{scopeFactory: scopeFactory}
}

func (floatingipReconcilerConstructor) GetName() string {
	return controllerName
}

const controllerName = "floatingip"

var (
	fieldOwner = orcstrings.GetSSAFieldOwner(controllerName)

	// FloatingIP depends on its external network
	networkDep = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Network](
		"spec.resource.networkRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil {
				return nil
			}
			return []string{string(resource.NetworkRef)}
		},
		finalizer, fieldOwner,
	)

	subnetDep = dependency.NewDeletionGuardDependency[*orcv1alpha1.FloatingIPList, *orcv1alpha1.Subnet](
		"spec.subnetRef",
		func(floatingip *orcv1alpha1.FloatingIP) []string {
			resource := floatingip.Spec.Resource
			if resource == nil {
				return nil
			}
			if resource.SubnetRef == nil {
				return nil
			}
			return []string{string(*resource.SubnetRef)}
		},
		finalizer, fieldOwner,
	)
)

// SetupWithManager sets up the controller with the Manager.
func (c floatingipReconcilerConstructor) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := mgr.GetLogger().WithValues("controller", controllerName)
	k8sClient := mgr.GetClient()

	networkHandler, err := networkDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	subnetHandler, err := subnetDep.WatchEventHandler(log, k8sClient)
	if err != nil {
		return err
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&orcv1alpha1.FloatingIP{}).
		Watches(&orcv1alpha1.Network{}, networkHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Network{})),
		).
		Watches(&orcv1alpha1.Subnet{}, subnetHandler,
			builder.WithPredicates(predicates.NewBecameAvailable(log, &orcv1alpha1.Subnet{})),
		)

	if err := errors.Join(
		networkDep.AddToManager(ctx, mgr),
		credentialsDependency.AddToManager(ctx, mgr),
		credentials.AddCredentialsWatch(log, k8sClient, builder, credentialsDependency),
	); err != nil {
		return err
	}

	r := reconciler.NewController(controllerName, k8sClient, c.scopeFactory, floatingipHelperFactory{}, floatingipStatusWriter{})
	return builder.Complete(&r)
}
