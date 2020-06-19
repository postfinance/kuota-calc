package calc

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment calculates the cpu/memory resources a single deployment needs. Replicas and the deployment
// strategy are taken into account.
func Deployment(deployment appsv1.Deployment) (*ResourceUsage, error) {
	var overhead float64

	replicas := deployment.Spec.Replicas
	strategy := deployment.Spec.Strategy

	switch strategy.Type {
	case appsv1.RecreateDeploymentStrategyType:
		// no overhead on recreate
		overhead = 1
	case appsv1.RollingUpdateDeploymentStrategyType:
		// As per https://pkg.go.dev/k8s.io/api/apps/v1?tab=doc#RollingUpdateDeployment absolute number is calculated
		// by rounding down.
		maxUnavailable, err := intstr.GetValueFromIntOrPercent(strategy.RollingUpdate.MaxUnavailable, int(*replicas), false)
		if err != nil {
			return nil, err
		}

		// As per https://pkg.go.dev/k8s.io/api/apps/v1?tab=doc#RollingUpdateDeployment absolute number is calculated
		// by rounding up.
		maxSurge, err := intstr.GetValueFromIntOrPercent(strategy.RollingUpdate.MaxSurge, int(*replicas), true)
		if err != nil {
			return nil, err
		}

		// podOverhead is the number of pods which can run more during a deployment
		podOverhead := maxSurge - maxUnavailable

		overhead = (float64(podOverhead) / float64(*replicas)) + 1
	default:
		return nil, fmt.Errorf("deployment strategy %s is not yet known", strategy.Type)
	}

	var (
		cpu    = new(resource.Quantity)
		memory = new(resource.Quantity)
	)

	// TODO: handle initContainers
	for _, container := range deployment.Spec.Template.Spec.Containers {
		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	resourceUsage := ResourceUsage{
		CPU:      float64(cpu.ScaledValue(resource.Kilo)*int64(*replicas)) * overhead,
		Memory:   float64(memory.ScaledValue(resource.Mega)*int64(*replicas)) * overhead,
		Overhead: overhead * 100,
	}

	return &resourceUsage, nil
}
