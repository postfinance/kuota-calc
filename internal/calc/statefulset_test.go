package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/deprecated/scheme"
)

var normalStatefulSet = `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  replicas: 2
  revisionHistoryLimit: 1
  selector:
    matchLabels:
      app: myapp
  updateStrategy:
    type: RollingUpdate
  serviceName: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - image: myapp
        imagePullPolicy: Always
        name: myapp
        resources:
          limits:
            cpu: "1"
            memory: 4Gi
          requests:
            cpu: 250m
            memory: 2Gi
      terminationGracePeriodSeconds: 30`

var noReplicasStatefulSet = `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  revisionHistoryLimit: 1
  selector:
    matchLabels:
      app: myapp
  updateStrategy:
    type: RollingUpdate
  serviceName: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - image: myapp
        imagePullPolicy: Always
        name: myapp
        resources:
          limits:
            cpu: "1"
            memory: 4Gi
          requests:
            cpu: 250m
            memory: 2Gi
      terminationGracePeriodSeconds: 30`

func TestStatefulSet(t *testing.T) {
	var tests = []struct {
		name        string
		statefulset appsv1.StatefulSet
		cpu         resource.Quantity
		memory      resource.Quantity
		replicas    int32
		maxReplicas int32
		strategy    appsv1.StatefulSetUpdateStrategyType
	}{
		{
			name:        "ok",
			statefulset: *statefulsetFromYaml(t, normalStatefulSet),
			cpu:         resource.MustParse("2"),
			memory:      resource.MustParse("8Gi"),
			replicas:    2,
			maxReplicas: 2,
			strategy:    appsv1.RollingUpdateStatefulSetStrategyType,
		},
		{
			name:        "no replicas",
			statefulset: *statefulsetFromYaml(t, noReplicasStatefulSet),
			cpu:         resource.MustParse("1"),
			memory:      resource.MustParse("4Gi"),
			replicas:    1,
			maxReplicas: 1,
			strategy:    appsv1.RollingUpdateStatefulSetStrategyType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			usage, err := StatefulSet(test.statefulset)
			r.NoError(err)
			r.NotEmpty(usage)

			r.Equalf(test.cpu.Value(), usage.CPU.Value(), "cpu value")
			r.Equalf(test.memory.Value(), usage.Memory.Value(), "memory value")
			r.Equalf(test.replicas, usage.Details.Replicas, "replicas")
			r.Equalf(test.maxReplicas, usage.Details.MaxReplicas, "maxReplicas")
			r.Equalf(string(test.strategy), usage.Details.Strategy, "strategy")
		})
	}
}

func statefulsetFromYaml(t *testing.T, statefulset string) *appsv1.StatefulSet {
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(statefulset), nil, nil)
	if err != nil {
		t.Error(err)
	}

	return object.(*appsv1.StatefulSet)
}
