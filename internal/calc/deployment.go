package calc

import (
	"fmt"
	"math"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment calculates the cpu/memory resources a single deployment needs. Replicas and the deployment
// strategy are taken into account.
func Deployment(deployment appsv1.Deployment) (*ResourceUsage, error) {
	var overhead float64

	replicas := deployment.Spec.Replicas

	if *replicas == 0 {
		return &ResourceUsage{
			CPU:    new(resource.Quantity),
			Memory: new(resource.Quantity),
			Details: Details{
				Version:  deployment.APIVersion,
				Kind:     deployment.Kind,
				Name:     deployment.Name,
				Replicas: *replicas,
			},
		}, nil
	}

	strategy := deployment.Spec.Strategy

	switch strategy.Type {
	case appsv1.RecreateDeploymentStrategyType:
		// no overhead on recreate
		overhead = 1
	case "":
		// RollingUpdate is the default an can be an empty string. If so, set the defaults
		// (https://pkg.go.dev/k8s.io/api/apps/v1?tab=doc#RollingUpdateDeployment) and continue calculation.
		defaults := intstr.FromString("25%")
		strategy = appsv1.DeploymentStrategy{
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: &defaults,
				MaxSurge:       &defaults,
			},
		}

		fallthrough
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
		return nil, fmt.Errorf("deployment: %s deployment strategy %q is unknown", deployment.Name, strategy.Type)
	}

	cpu, memory := podResources(&deployment.Spec.Template.Spec)

	mem := float64(memory.Value()) * float64(*replicas) * overhead
	memory.Set(int64(math.Round(mem)))

	cpu.Set(int64(math.Round(float64(cpu.Value()) * float64(*replicas) * overhead)))

	resourceUsage := ResourceUsage{
		CPU:    cpu,
		Memory: memory,
		Details: Details{
			Version:  deployment.APIVersion,
			Kind:     deployment.Kind,
			Name:     deployment.Name,
			Replicas: *replicas,
		},
	}

	return &resourceUsage, nil
}
