package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/deprecated/scheme"
)

func TestDeployment(t *testing.T) {
	var tests = []struct {
		name        string
		deployment  appsv1.Deployment
		cpu         resource.Quantity
		memory      resource.Quantity
		replicas    int32
		maxReplicas int32
		strategy    appsv1.DeploymentStrategyType
	}{
		{
			name:        "normal deployment",
			deployment:  *deploymentFromYaml(t, normalDeployment),
			cpu:         resource.MustParse("5500m"),
			memory:      resource.MustParse("44Gi"),
			replicas:    10,
			maxReplicas: 11,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
		{
			name:        "deployment without strategy",
			deployment:  *deploymentFromYaml(t, deploymentWithoutStrategy),
			cpu:         resource.MustParse("11"),
			memory:      resource.MustParse("44Gi"),
			replicas:    10,
			maxReplicas: 11,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
		{
			name:        "deployment with absolute unavailable/surge values",
			deployment:  *deploymentFromYaml(t, deploymentWithAbsoluteValues),
			cpu:         resource.MustParse("12"),
			memory:      resource.MustParse("48Gi"),
			replicas:    10,
			maxReplicas: 12,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
		{
			name:        "zero replica deployment",
			deployment:  *deploymentFromYaml(t, zeroReplicaDeployment),
			cpu:         resource.MustParse("0"),
			memory:      resource.MustParse("0"),
			replicas:    0,
			maxReplicas: 0,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
		{
			name:        "recreate deployment",
			deployment:  *deploymentFromYaml(t, recrateDeployment),
			cpu:         resource.MustParse("10"),
			memory:      resource.MustParse("40Gi"),
			replicas:    10,
			maxReplicas: 10,
			strategy:    appsv1.RecreateDeploymentStrategyType,
		},
		{
			name:        "deployment without max unavailable/surge values",
			deployment:  *deploymentFromYaml(t, deploymentWithoutValues),
			cpu:         resource.MustParse("11"),
			memory:      resource.MustParse("44Gi"),
			replicas:    10,
			maxReplicas: 11,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
		{
			name:        "deployment with init container(s)",
			deployment:  *deploymentFromYaml(t, initContainerDeployment),
			cpu:         resource.MustParse("4400m"),
			memory:      resource.MustParse("17184Mi"),
			replicas:    3,
			maxReplicas: 4,
			strategy:    appsv1.RollingUpdateDeploymentStrategyType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			usage, err := Deployment(test.deployment)
			r.NoError(err)
			r.NotEmpty(usage)

			r.Equalf(test.cpu.MilliValue(), usage.CPU.MilliValue(), "cpu value")
			r.Equal(0, test.memory.Cmp(*usage.Memory), "memory value %d != %d", test.memory.Value(), usage.Memory.Value())
			r.Equal(test.replicas, usage.Details.Replicas, "replicas")
			r.Equal(string(test.strategy), usage.Details.Strategy, "strategy")
			r.Equal(test.maxReplicas, usage.Details.MaxReplicas, "maxReplicas")
		})
	}
}

func deploymentFromYaml(t *testing.T, deployment string) *appsv1.Deployment {
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(deployment), nil, nil)
	if err != nil {
		t.Error(err)
	}

	return object.(*appsv1.Deployment)
}
