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
		name             string
		deployment       appsv1.Deployment
		expectedCPU      resource.Quantity
		expectedMemory   resource.Quantity
		expectedReplicas int32
	}{
		{
			name:             "normal deployment",
			deployment:       *deploymentFromYaml(t, normalDeployment),
			expectedCPU:      resource.MustParse("11"),
			expectedMemory:   resource.MustParse("44Gi"),
			expectedReplicas: 10,
		},
		{
			name:             "deployment without strategy",
			deployment:       *deploymentFromYaml(t, deploymentWithoutStrategy),
			expectedCPU:      resource.MustParse("11"),
			expectedMemory:   resource.MustParse("44Gi"),
			expectedReplicas: 10,
		},
		{
			name:             "deployment with absolute unavailable/surge values",
			deployment:       *deploymentFromYaml(t, deploymentWithAbsoluteValues),
			expectedCPU:      resource.MustParse("12"),
			expectedMemory:   resource.MustParse("48Gi"),
			expectedReplicas: 10,
		},
		{
			name:             "zero replica deployment",
			deployment:       *deploymentFromYaml(t, zeroReplicaDeployment),
			expectedCPU:      resource.MustParse("0"),
			expectedMemory:   resource.MustParse("0"),
			expectedReplicas: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			usage, err := Deployment(test.deployment)
			r.NoError(err)
			r.NotEmpty(usage)

			r.Equalf(test.expectedCPU.Value(), usage.CPU.Value(), "cpu value")
			r.Equalf(test.expectedMemory.Value(), usage.Memory.Value(), "memory value")
			r.Equal(test.expectedReplicas, usage.Details.Replicas, "replicas")
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
