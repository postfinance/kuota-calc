// Package calc provides function to calculate resource quotas for different k8s resources.
package calc

import (
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/deprecated/scheme"
)

var (
	// ErrResourceNotSupported is returned if a k8s resource is not supported by kuota-calc.
	ErrResourceNotSupported = errors.New("resource not supported")
)

// CalculationError is an error implementation that includes a k8s Kind/Version.
type CalculationError struct {
	Version string
	Kind    string
	err     error
}

func (cErr CalculationError) Error() string {
	return fmt.Sprintf("calculating %s/%s resource usage: %s",
		cErr.Version,
		cErr.Kind,
		cErr.err,
	)
}

// Unwrap implements the errors.Unwrap interface.
func (cErr CalculationError) Unwrap() error {
	return cErr.err
}

// ResourceUsage summarizes the usage of compute resources for a k8s resource.
type ResourceUsage struct {
	CPU     *resource.Quantity
	Memory  *resource.Quantity
	Details Details
}

// Details contains a few details of a k8s resource, which are needed to generate a detailed resource
// usage report.
type Details struct {
	Version     string
	Kind        string
	Name        string
	Strategy    string
	Replicas    int32
	MaxReplicas int32
}

func podResources(podSpec *v1.PodSpec) (cpu, memory *resource.Quantity) {
	cpu = new(resource.Quantity)
	memory = new(resource.Quantity)

	for _, container := range podSpec.Containers {
		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	for _, container := range podSpec.InitContainers {
		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	return
}

// ResourceQuotaFromYaml decodes a single yaml document into a k8s object. Then performs a type assertion
// on the object and calculates the resource needs of it.
// Currently supported:
// * apps/v1 - Deployment
// * apps/v1 - StatefulSet
func ResourceQuotaFromYaml(yamlData []byte) (*ResourceUsage, error) {
	object, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(yamlData, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("decoding yaml data: %w", err)
	}

	switch obj := object.(type) {
	case *appsv1.Deployment:
		usage, err := deployment(*obj)
		if err != nil {
			return nil, CalculationError{
				Version: gvk.Version,
				Kind:    gvk.Kind,
				err:     err,
			}
		}

		return usage, nil
	case *appsv1.StatefulSet:
		return statefulSet(*obj), nil
	default:
		return nil, CalculationError{
			Version: gvk.Version,
			Kind:    gvk.Kind,
			err:     ErrResourceNotSupported,
		}
	}
}
