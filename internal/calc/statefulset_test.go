package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/deprecated/scheme"
)

var testStatefulSet = `
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
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(testStatefulSet), nil, nil)
	if err != nil {
		t.Error(err)
	}
	statefulset := object.(*appsv1.StatefulSet)

	var tests = []struct {
		name           string
		statefulset    appsv1.StatefulSet
		expectedCPU    resource.Quantity
		expectedMemory resource.Quantity
	}{
		{
			name:           "ok",
			statefulset:    *statefulset,
			expectedCPU:    resource.MustParse("2"),
			expectedMemory: resource.MustParse("8Gi"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			usage, err := StatefulSet(test.statefulset)
			r.NoError(err)
			r.NotEmpty(usage)

			r.Equalf(test.expectedCPU.Value(), usage.CPU.Value(), "cpu value")
			r.Equalf(test.expectedMemory.Value(), usage.Memory.Value(), "memory value")
		})
	}
}
