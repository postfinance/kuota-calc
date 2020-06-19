![ci](https://github.com/postfinance/kuota-calc/workflows/ci/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/postfinance/kuota-calc)
[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/kuota-calc)](https://goreportcard.com/report/github.com/postfinance/kuota-calc)
![License](https://img.shields.io/github/license/postfinance/kuota-calc)

**This is still a work-in-progress**

# kuota-calc
Simple utility to calculate the resource quota needed for your deployment. kuota-calc takes the
deployment strategy, replicas and all containers into account.

## Example

Deployment:
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  replicas: 2
  selector:
    matchLabels:
      app: myapp
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
        labels:
        app: myapp
    spec:
      containers:
      - image: myapp:0.1.0
        imagePullPolicy: Always
        name: myapp
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 250m
            memory: 64Mi
```

```bash
$ cat deployment.yaml | kuota-calc -detailed
Version    Kind          Name     Replicas    CPU    Memory
apps/v1    Deployment    myapp    2           2      256Mi

Total
CPU: 2
Memory: 256Mi
```
