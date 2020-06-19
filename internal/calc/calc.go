// Package calc provides function to calculate resource quotas for different k8s resources.
package calc

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ResourceUsage summarizes the usage of compute resources for a k8s resource.
type ResourceUsage struct {
	CPU      *resource.Quantity
	Memory   *resource.Quantity
	Overhead float64
}

func podResources(podSpec *v1.PodSpec) (*resource.Quantity, *resource.Quantity) {
	var (
		cpu    = new(resource.Quantity)
		memory = new(resource.Quantity)
	)

	for _, container := range podSpec.Containers {
		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	for _, container := range podSpec.InitContainers {
		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	return cpu, memory
}
