package calc

import (
	"math"

	appsv1 "k8s.io/api/apps/v1"
)

// StatefulSet calculates the cpu/memory resources a single statefulset needs. Replicas are taken into account.
func StatefulSet(s appsv1.StatefulSet) (*ResourceUsage, error) {
	var (
		replicas int32
	)

	// https://github.com/kubernetes/api/blob/v0.18.4/apps/v1/types.go#L117
	if s.Spec.Replicas != nil {
		replicas = *s.Spec.Replicas
	} else {
		replicas = 1
	}

	cpu, memory := podResources(&s.Spec.Template.Spec)

	mem := float64(memory.Value()) * float64(replicas)
	memory.Set(int64(math.Round(mem)))

	cpu.Set(int64(math.Round(float64(cpu.Value()) * float64(replicas))))

	resourceUsage := ResourceUsage{
		CPU:    cpu,
		Memory: memory,
		Details: Details{
			Version:     s.APIVersion,
			Kind:        s.Kind,
			Name:        s.Name,
			Replicas:    replicas,
			Strategy:    string(s.Spec.UpdateStrategy.Type),
			MaxReplicas: replicas,
		},
	}

	return &resourceUsage, nil
}
