# Connectors

- **TODO**: Fill this

## Scale Alert Emitters

- [prometheus/alertmanager](https://prometheus.io/docs/alerting/alertmanager/)

## Scaling Targets

Sokar is able to scale several types of entities. These are:

- Nomad Job
- Nomad Nodes on AWS
- AWS instances
- Kubernetes pods (not implemented yet)

For each of them sokar has to be configured to run in the according mode and depending on the scaling target type additional configuration parameters (e.g. credentials or a server address) have to be specified.

How to run sokar in these different modes, please have a look at [ScalerModes.md](ScalerModes.md).
