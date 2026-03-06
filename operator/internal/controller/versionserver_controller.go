/*
Copyright 2026.

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

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	metalk8ssolutionexamplescalitycomv1alpha1 "github.com/scality/metalk8s-solution-example/operator/api/v1alpha1"
	"github.com/scality/metalk8s-solution-example/operator/internal/utils"
	opConfig "github.com/scality/metalk8s/go/solution-operator-lib/pkg/config"
)

const (
	versionServerControllerName = "versionserver"
	versionServerAppName        = "versionserver"
)

// VersionServerReconciler reconciles a VersionServer object
type VersionServerReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	OperatorConfig *opConfig.OperatorConfig
}

// +kubebuilder:rbac:groups=metalk8s-solution-example.scality.com,resources=versionservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metalk8s-solution-example.scality.com,resources=versionservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metalk8s-solution-example.scality.com,resources=versionservers/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VersionServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *VersionServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	logger.Info("Reconciling VersionServer: START")
	defer logger.Info("Reconciling VersionServer: STOP")

	// Fetch the VersionServer instance
	instance := &metalk8ssolutionexamplescalitycomv1alpha1.VersionServer{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// --- Deployment ---
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		err := utils.MutateBaseServerDeployment(
			deployment,
			instance,
			r.Scheme,
			versionServerControllerName,
			versionServerAppName,
		)
		if err != nil {
			return err
		}

		err = utils.MutateBaseServerPodSpec(
			&deployment.Spec.Template.Spec,
			versionServerAppName,
			instance.Spec.Version,
			r.OperatorConfig.Repositories,
			[]string{"--version", instance.Spec.Version},
		)
		if err != nil {
			return err
		}

		deployment.Spec.Replicas = &instance.Spec.Replicas
		return nil
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	if op != controllerutil.OperationResultNone {
		logger.Info("Deployment for VersionServer reconciled", "operation", op)
	}

	// --- Service ---
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		return utils.MutateBaseServerService(
			service,
			instance,
			r.Scheme,
			versionServerControllerName,
			versionServerAppName,
		)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	if op != controllerutil.OperationResultNone {
		logger.Info("Service for VersionServer reconciled", "operation", op)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VersionServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalk8ssolutionexamplescalitycomv1alpha1.VersionServer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named(versionServerControllerName).
		Complete(r)
}
