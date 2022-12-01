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
	"fmt"

	kapps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/x893675/demo-crd/apis/apps/v1"
)

// SimpleDeploymentReconciler reconciles a SimpleDeployment object
type SimpleDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=simpledeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=simpledeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.x893675.github.io,resources=simpledeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SimpleDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *SimpleDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx).WithValues("simpleDeployment", req.NamespacedName)

	var simpleDeployment appsv1.SimpleDeployment
	if err := r.Get(ctx, req.NamespacedName, &simpleDeployment); err != nil {
		logs.Error(err, "unable to fetch SimpleDeployment")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// +kubebuilder:docs-gen:collapse=Begin the Reconcile

	/*
		Build the deployment that we want to see exist within the cluster
	*/

	foundDeployment := &kapps.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: simpleDeployment.Name, Namespace: simpleDeployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logs.V(1).Info("Creating Deployment", "deployment", simpleDeployment.Name)
		deployment := initDeployment(simpleDeployment.Name, simpleDeployment.Namespace)
		if err := controllerutil.SetControllerReference(&simpleDeployment, deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		err = r.Create(ctx, deployment)
	} else if err == nil {
		if foundDeployment.Spec.Replicas != simpleDeployment.Spec.Replicas {
			foundDeployment.Spec.Replicas = simpleDeployment.Spec.Replicas
			logs.V(1).Info("Updating Deployment", "deployment", simpleDeployment.Name)
			err = r.Update(ctx, foundDeployment)
		}
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.SimpleDeployment{}).
		Owns(&kapps.Deployment{}).
		Complete(r)
}

func initDeployment(name, ns string) *kapps.Deployment {
	return &kapps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"app": fmt.Sprintf("d-%s", name),
			},
		},
		Spec: kapps.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": fmt.Sprintf("d-%s", name),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": fmt.Sprintf("d-%s", name),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}
