package utils

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/scality/metalk8s-solution-example/operator/internal"
)

const (
	// Common labels
	labelManagedByName  = "app.kubernetes.io/managed-by"
	labelManagedByValue = "example-solution-operator"

	labelPartOfName  = "app.kubernetes.io/part-of"
	labelPartOfValue = "metalk8s-solution-example"

	labelComponentName = "app.kubernetes.io/component"
	labelAppName       = "app.kubernetes.io/name"
	labelInstanceName  = "app.kubernetes.io/instance"

	labelVersionName = "app.kubernetes.io/version"
)

var (
	stdLabels = map[string]string{
		labelManagedByName: labelManagedByValue,
		labelPartOfName:    labelPartOfValue,
		labelVersionName:   internal.Version,
	}

	// List of labels used as Selector for Deployment pods
	selectorLabels = []string{
		labelAppName,
		labelInstanceName,
	}
)

// getCommonLabels returns a map of common labels
func getCommonLabels(
	controllerName, appName, instanceName string,
) map[string]string {
	labels := make(map[string]string)
	for k, v := range stdLabels {
		labels[k] = v
	}
	labels[labelComponentName] = controllerName
	labels[labelAppName] = appName
	labels[labelInstanceName] = instanceName
	return labels
}

// getSelectorLabels returns a map of selector labels based on all labels in the map
func getSelectorLabels(labels map[string]string) client.MatchingLabels {
	selector := make(client.MatchingLabels)
	for _, label := range selectorLabels {
		selector[label] = labels[label]
	}

	return selector
}

// StdMutate mutates an object with common labels and selector labels
func StdMutate(object metav1.Object, instance metav1.Object, scheme *runtime.Scheme, controllerName, appName string) error {
	commonLabels := getCommonLabels(controllerName, appName, instance.GetName())
	UpdateLabels(object, commonLabels)

	switch object := object.(type) {
	case *appsv1.Deployment:
		err := stdMutateDeployment(object, commonLabels)
		if err != nil {
			return err
		}
	}

	return controllerutil.SetControllerReference(instance, object, scheme)
}

// stdMutateDeployment mutates a Deployment with common labels and selector labels
// nolint:unparam // This function may return error in the future
func stdMutateDeployment(deployment *appsv1.Deployment, commonLabels map[string]string) error {
	UpdateLabels(&deployment.Spec.Template.ObjectMeta, commonLabels)

	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: getSelectorLabels(commonLabels),
	}
	return nil
}

// Update labels on an Object
func UpdateLabels(object metav1.Object, labels map[string]string) {
	current := object.GetLabels()
	if current == nil {
		current = make(map[string]string)
	}

	for k, v := range labels {
		current[k] = v
	}

	object.SetLabels(current)
}
