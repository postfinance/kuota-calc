![ci](https://github.com/postfinance/kuota-calc/workflows/ci/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/postfinance/kuota-calc)
[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/kuota-calc)](https://goreportcard.com/report/github.com/postfinance/kuota-calc)
![License](https://img.shields.io/github/license/postfinance/kuota-calc)

# kuota-calc
Simple utility to calculate the maximum needed resource quota for deployment(s). kuota-calc takes the
deployment strategy, replicas and all containers into account, see [supported-resources](https://github.com/postfinance/kuota-calc#supported-k8s-resources) for a list of kubernetes resources which are currently supported by kuota-calc.

## Motivation
In shared environments such as kubernetes it is always a good idea to isolate/constrain different workloads to prevent them from infering each other. Kubernetes provides [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/) to limit compute, storage and object resources of namespaces.

Calculating the needed compute resources can be a bit challenging (especially with large and complex deployments) because we must respect certain settings/defaults like the deployment strategy, number of replicas and so on. This is where kuota-calc can help you, it calculates the maximum needed resource quota in order to be able to start a deployment of all resources at the same time by respecting deployment strategies, replicas and so on.

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

## supported k8s resources
**kuota-calc is still a work-in progress**, there are plans to support more k8s resources (see [#5](https://github.com/postfinance/kuota-calc/issues/5) for more info). 

Currently supported:

- apps/v1 Deployment
- apps/v1 StatefulSet
- apps/v1 DaemonSet
- batch/v1 CronJob
- batch/v1 Job
- v1 Pod
