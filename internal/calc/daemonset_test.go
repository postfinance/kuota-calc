package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDaemonSet(t *testing.T) {
	var tests = []struct {
		name      string
		daemonset string
		cpu       resource.Quantity
		memory      resource.Quantity
		replicas    int32
		maxReplicas int32
		strategy    appsv1.StatefulSetUpdateStrategyType
	}{
		{
			name:      "ok",
			daemonset: normalDaemonSet,
			replicas: 1,
			maxReplicas: 1,
			cpu:       resource.MustParse("2"),
			memory:    resource.MustParse("2Gi"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				r := require.New(t)

				usage, err := ResourceQuotaFromYaml([]byte(test.daemonset))
				r.NoError(err)
				r.NotEmpty(usage)

				r.Equalf(test.cpu.Value(), usage.CPU.Value(), "cpu value")
				r.Equalf(test.memory.Value(), usage.Memory.Value(), "memory value")
				r.Equalf(test.replicas, usage.Details.Replicas, "replicas")
				r.Equalf(test.maxReplicas, usage.Details.MaxReplicas, "maxReplicas")
				r.Equalf(string(test.strategy), usage.Details.Strategy, "strategy")
			},
		)
	}
}
