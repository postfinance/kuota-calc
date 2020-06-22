![ci](https://github.com/postfinance/kuota-calc/workflows/ci/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/postfinance/kuota-calc)
[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/kuota-calc)](https://goreportcard.com/report/github.com/postfinance/kuota-calc)
![License](https://img.shields.io/github/license/postfinance/kuota-calc)

**This is still a work-in-progress**

# kuota-calc
Simple utility to calculate the resource quota needed for your deployment. kuota-calc takes the
deployment strategy, replicas and all containers into account.

## Example
```bash
$ cat examples/deployment.yaml | kuota-calc -detailed
Version    Kind           Name     Replicas    Strategy         MaxReplicas    CPU      Memory
apps/v1    Deployment     myapp    10          RollingUpdate    11             5500m    2816Mi
apps/v1    StatefulSet    myapp    3           RollingUpdate    3              3        12Gi

Total
CPU: 8500m
Memory: 15104Mi
```
