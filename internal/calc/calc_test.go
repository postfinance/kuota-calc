package calc

var normalDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: normal
  name: normal
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: normal
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: normal
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: normal
          resources:
            limits:
              cpu: '500m'
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

var initContainerDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: normal
  name: normal
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: normal
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: normal
    spec:
      initContainers:
        - image: myinit:v1.0.7
          name: myinit
          resources:
            limits:
              cpu: '100m'
              memory: '200Mi'
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: normal
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

var recrateDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: normal
  name: normal
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: normal
  strategy:
    type: Recreate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: normal
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: normal
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

var zeroReplicaDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: zero
  name: zero
spec:
  progressDeadlineSeconds: 600
  replicas: 0
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: zero
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: zero
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: zero
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

var deploymentWithoutStrategy = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: strategy
  name: strategy
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: strategy
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: strategy
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: strategy
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

var deploymentWithAbsoluteValues = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: values
  name: values
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: values
  strategy:
    rollingUpdate:
      maxSurge: 2
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: values
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: values
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

var deploymentWithoutValues = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: values
  name: values
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: values
  strategy:
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: values
    spec:
      containers:
        - image: myapp:v1.0.7
          command:
            - sleep
            - infinity
          imagePullPolicy: IfNotPresent
          name: values
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
