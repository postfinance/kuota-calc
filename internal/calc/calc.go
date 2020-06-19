// Package calc provides function to calculate resource quotas for different k8s resources.
package calc

import "k8s.io/apimachinery/pkg/api/resource"

// ResourceUsage summarizes the usage of compute resources for a k8s resource.
type ResourceUsage struct {
	CPU      *resource.Quantity
	Memory   *resource.Quantity
	Overhead float64
}
