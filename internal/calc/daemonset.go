package calc

import (
	appsv1 "k8s.io/api/apps/v1"
)

func daemonSet(dSet appsv1.DaemonSet) *ResourceUsage {

	cpu, memory := podResources(&dSet.Spec.Template.Spec)

	resourceUsage := ResourceUsage{
		CPU:    cpu,
		Memory: memory,
		Details: Details{
			Version:     dSet.APIVersion,
			Kind:        dSet.Kind,
			Name:        dSet.Name,
			Strategy:    "",
			Replicas:    0,
			MaxReplicas: 0,
		},
	}

	return &resourceUsage
}
