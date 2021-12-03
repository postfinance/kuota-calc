// Package calc provides function to calculate resource quotas for different k8s resources.
package calc

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
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
	return fmt.Sprintf(
		"calculating %s/%s resource usage: %s",
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

	for i := range podSpec.Containers {
		container := podSpec.Containers[i]

		cpu.Add(*container.Resources.Limits.Cpu())
		memory.Add(*container.Resources.Limits.Memory())
	}

	for i := range podSpec.InitContainers {
		container := podSpec.InitContainers[i]

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
// * apps/v1 - DaemonSet
// * v1 - Pod
type banana struct {
	schemeName string
	gvk        schema.GroupVersionKind
	target     runtime.GroupVersioner
	t          reflect.Type
}

func (b *banana) Error() string {
	return ""
}

func ResourceQuotaFromYaml(yamlData []byte) (*ResourceUsage, error) {

	var version string
	var kind string

	object, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(yamlData, nil, nil)

	if err != nil {
		// when the kind is not found, I just warn and skip
		if runtime.IsNotRegisteredError(err) {
			log.Warn().Msg(err.Error())
			unknown := runtime.Unknown{Raw: yamlData}

			if _, gvk1, err := scheme.Codecs.UniversalDeserializer().Decode(yamlData, nil, &unknown); err == nil {
				kind = gvk1.Kind
				version = gvk1.Version
			}
		} else {
			return nil, fmt.Errorf("decoding yaml data: %w", err)
		}
	} else {
		kind = gvk.Kind
		version = gvk.Version
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
	case *appsv1.DaemonSet:
		return daemonSet(*obj), nil
	case *batchV1.Job:
		return job(*obj), nil
	case *batchV1.CronJob:
		return cronjob(*obj), nil
	case *v1.Pod:
		return pod(*obj), nil
	default:
		return nil, CalculationError{
			Version: version,
			Kind:    kind,
			err:     ErrResourceNotSupported,
		}
	}
}
