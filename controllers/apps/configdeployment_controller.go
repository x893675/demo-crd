/*
Copyright 2022.

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

package apps

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "github.com/x893675/demo-crd/apis/apps/v1"
	kapps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

/*
Determine the path of the field in the ConfigDeployment CRD that we wish to use as the "object reference".
This will be used in both the indexing and watching.
*/
const (
	configMapField = ".spec.configMap"
)

// ConfigDeploymentReconciler reconciles a ConfigDeployment object
type ConfigDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=configdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=configdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=configdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ConfigDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx)

	var configDeployment appsv1.ConfigDeployment
	if err := r.Get(ctx, req.NamespacedName, &configDeployment); err != nil {
		logs.Error(err, "unable to fetch ConfigDeployment")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// your logic here

	var configMapVersion string
	if configDeployment.Spec.ConfigMap != "" {
		configMapName := configDeployment.Spec.ConfigMap
		foundConfigMap := &corev1.ConfigMap{}
		err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: configDeployment.Namespace}, foundConfigMap)
		if err != nil {
			// If a configMap name is provided, then it must exist
			// You will likely want to create an Event for the user to understand why their reconcile is failing.
			return ctrl.Result{}, err
		}

		// Hash the data in some way, or just use the version of the Object
		configMapVersion = foundConfigMap.ResourceVersion
	}

	logs.Info("configMapVersion is %s", configMapVersion)
	// Logic here to add the configMapVersion as an annotation on your Deployment Pods.

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	/*
		The `configMap` field must be indexed by the manager, so that we will be able to lookup `ConfigDeployments` by a referenced `ConfigMap` name.
		This will allow for quickly answer the question:
		- If ConfigMap _x_ is updated, which ConfigDeployments are affected?
	*/

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.ConfigDeployment{}, configMapField, func(rawObj client.Object) []string {
		// Extract the ConfigMap name from the ConfigDeployment Spec, if one is provided
		configDeployment := rawObj.(*appsv1.ConfigDeployment)
		if configDeployment.Spec.ConfigMap == "" {
			return nil
		}
		return []string{configDeployment.Spec.ConfigMap}
	}); err != nil {
		return err
	}

	/*
		As explained in the CronJob tutorial, the controller will first register the Type that it manages, as well as the types of subresources that it controls.
		Since we also want to watch ConfigMaps that are not controlled or managed by the controller, we will need to use the `Watches()` functionality as well.
		The `Watches()` function is a controller-runtime API that takes:
		- A Kind (i.e. `ConfigMap`)
		- A mapping function that converts a `ConfigMap` object to a list of reconcile requests for `ConfigDeployments`.
		We have separated this out into a separate function.
		- A list of options for watching the `ConfigMaps`
		  - In our case, we only want the watch to be triggered when the ResourceVersion of the ConfigMap is changed.
	*/

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ConfigDeployment{}).
		Owns(&kapps.Deployment{}).
		Watches(
			&source.Kind{Type: &corev1.ConfigMap{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

/*
Because we have already created an index on the `configMap` reference field, this mapping function is quite straight forward.
We first need to list out all `ConfigDeployments` that use `ConfigMap` given in the mapping function.
This is done by merely submitting a List request using our indexed field as the field selector.
When the list of `ConfigDeployments` that reference the `ConfigMap` is found,
we just need to loop through the list and create a reconcile request for each one.
If an error occurs fetching the list, or no `ConfigDeployments` are found, then no reconcile requests will be returned.
*/
func (r *ConfigDeploymentReconciler) findObjectsForConfigMap(configMap client.Object) []reconcile.Request {
	attachedConfigDeployments := &appsv1.ConfigDeploymentList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapField, configMap.GetName()),
		Namespace:     configMap.GetNamespace(),
	}
	err := r.List(context.TODO(), attachedConfigDeployments, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(attachedConfigDeployments.Items))
	for i, item := range attachedConfigDeployments.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}
	return requests
}
