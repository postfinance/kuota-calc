package calc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/deprecated/scheme"
)

var testDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: rhel
  name: rhel
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: rhel
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: rhel
    spec:
      containers:
        - image: linux-docker-local.repo.pnet.ch/pf/rhel:7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: rhel
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
	object, _, _ := scheme.Codecs.UniversalDeserializer().Decode([]byte(testDeployment), nil, nil)
	deployment := object.(*appsv1.Deployment)

	var tests = []struct {
		name          string
		deployment    appsv1.Deployment
		expectedUsage ResourceUsage
	}{
		{
			name:       "ok",
			deployment: *deployment,
			expectedUsage: ResourceUsage{
				CPU:      11,
				Memory:   45056,
				Overhead: 110,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			usage, err := Deployment(test.deployment)
			r.NoError(err)
			r.NotEmpty(usage)

			r.Equal(test.expectedUsage.CPU, math.Round(usage.CPU))
			r.Equal(test.expectedUsage.Memory, math.Round(usage.Memory))
			r.Equal(test.expectedUsage.Overhead, math.Round(usage.Overhead))
		})
	}
}
