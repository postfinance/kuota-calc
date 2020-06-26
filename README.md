![ci](https://github.com/postfinance/kuota-calc/workflows/ci/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/postfinance/kuota-calc)
[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/kuota-calc)](https://goreportcard.com/report/github.com/postfinance/kuota-calc)
![License](https://img.shields.io/github/license/postfinance/kuota-calc)

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

## Installation
Pre-compiled statically linked binaries are available on the [releases page](https://github.com/postfinance/kuota-calc/releases).

kuota-calc can either be used as a kubectl plugin or invoked directly. If you intend to use kuota-calc as
a kubectl plugin, simply place the binary anywhere in `$PATH` named `kubectl-kuota_calc` with execute permissions.
For further information, see the offical documentation on kubectl plugins [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).
