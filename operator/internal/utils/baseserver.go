package utils

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	opConfig "github.com/scality/metalk8s/go/solution-operator-lib/pkg/config"
)

const (
	// ContainerHTTPPort is the HTTP port number used by `base-server`
	ContainerHTTPPort = 8080

	HEALTHZEndpoint = "/healthz"
)

// buildImageName builds a complete image name based on the version provided
// for the `base-server` component, which is the only one deployed for now
// Here `version` is both the image and the Solution versions
func buildImageName(version string, repositories map[string][]opConfig.Repository) (string, error) {
	var prefix string
	imageName := "base-server:" + version

	for solution_version, repositories := range repositories {
		if solution_version == version {
			for _, repository := range repositories {
				for _, image := range repository.Images {
					if image == imageName {
						prefix = repository.Endpoint
					}
				}
			}
		}
	}

	if prefix == "" {
		return "", fmt.Errorf(
			"unable to find image %s in repositories configuration",
			imageName,
		)
	}

	return fmt.Sprintf("%s/%s", prefix, imageName), nil
}

// MutateBaseServerDeployment mutates a Deployment with the base-server content
func MutateBaseServerDeployment(
	deployment *appsv1.Deployment,
	instance metav1.Object,
	scheme *runtime.Scheme,
	controllerName, appName string,
) error {
	err := StdMutate(deployment, instance, scheme, controllerName, appName)
	if err != nil {
		return err
	}

	maxSurge := intstr.FromInt(0)
	maxUnavailable := intstr.FromInt(1)

	deployment.Spec.Strategy = appsv1.DeploymentStrategy{
		Type: appsv1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       &maxSurge,
			MaxUnavailable: &maxUnavailable,
		},
	}

	return nil
}

// MutateBaseServerPodSpec mutates a PodSpec with the base-server content
func MutateBaseServerPodSpec(
	pod *corev1.PodSpec,
	appName string,
	imageVersion string,
	repositories map[string][]opConfig.Repository,
	cmdArgs []string,
) error {
	if pod.Containers == nil {
		pod.Containers = make([]corev1.Container, 1)
	}

	image, err := buildImageName(imageVersion, repositories)
	if err != nil {
		return err
	}
	pod.Containers[0].Name = fmt.Sprintf("%s-base-server", appName)
	pod.Containers[0].Image = image
	pod.Containers[0].ImagePullPolicy = corev1.PullIfNotPresent
	pod.Containers[0].Command = append([]string{"python3", "/app/server.py"}, cmdArgs...)
	pod.Containers[0].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   HEALTHZEndpoint,
				Port:   intstr.FromInt(ContainerHTTPPort),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		FailureThreshold:    8,
		InitialDelaySeconds: 10,
		TimeoutSeconds:      3,
	}
	pod.Containers[0].Ports = []corev1.ContainerPort{
		{
			ContainerPort: ContainerHTTPPort,
			Protocol:      corev1.ProtocolTCP,
			Name:          "http",
		},
	}

	return nil
}

// MutateBaseServerService mutates a Service with the base-server content
func MutateBaseServerService(
	service *corev1.Service,
	instance metav1.Object,
	scheme *runtime.Scheme,
	controllerName, appName string,
) error {
	err := StdMutate(service, instance, scheme, controllerName, appName)
	if err != nil {
		return err
	}

	service.Spec.Ports = []corev1.ServicePort{
		{
			Port:       ContainerHTTPPort,
			Protocol:   corev1.ProtocolTCP,
			Name:       "http",
			TargetPort: intstr.FromInt(ContainerHTTPPort),
		},
	}
	service.Spec.Selector = getSelectorLabels(service.GetLabels())
	return nil
}
