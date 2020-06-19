// Package calc provides function to calculate resource quotas for different k8s resources.
package calc

// ResourceUsage summarizes the usage of compute resources for a k8s resource.
type ResourceUsage struct {
	CPU      float64
	Memory   float64
	Overhead float64
}
