package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestJob(t *testing.T) {
	var tests = []struct {
		name        string
		job         string
		cpu         resource.Quantity
		memory      resource.Quantity
		replicas    int32
		maxReplicas int32
		strategy    string
	}{
		{
			name:   "ok",
			job:    normalJob,
			cpu:    resource.MustParse("1"),
			memory: resource.MustParse("4Gi"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				r := require.New(t)

				usage, err := ResourceQuotaFromYaml([]byte(test.job))
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
