package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/deprecated/scheme"
)

var testDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: myapp
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: myapp
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: myapp
          resources:
            limits:
              cpu: '1'
              memory: 4Gi
            requests:
              cpu: '250m'
              memory: 2Gi
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30`

func TestDeployment(t *testing.T) {
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(testDeployment), nil, nil)
	if err != nil {
		t.Error(err)
	}
	deployment := object.(*appsv1.Deployment)

	var tests = []struct {
		name           string
		deployment     appsv1.Deployment
		expectedCPU    resource.Quantity
		expectedMemory resource.Quantity
	}{
		{
			name:           "ok",
			deployment:     *deployment,
			expectedCPU:    resource.MustParse("11"),
			expectedMemory: resource.MustParse("44Gi"),
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
		})
	}
}
