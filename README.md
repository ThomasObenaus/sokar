# Sokar

[![Go Report Card](https://goreportcard.com/badge/github.com/ThomasObenaus/sokar)](https://goreportcard.com/report/github.com/ThomasObenaus/sokar) [![Maintainability](https://api.codeclimate.com/v1/badges/56824372d45781170a68/maintainability)](https://codeclimate.com/github/ThomasObenaus/sokar/maintainability) [![Coverage Status](https://coveralls.io/repos/github/ThomasObenaus/sokar/badge.svg?branch=master)](https://coveralls.io/github/ThomasObenaus/sokar?branch=master) [![FOSSA Status](https://app.fossa.com/api/projects/custom%2B12599%2Fgit%40github.com%3AThomasObenaus%2Fsokar.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B12599%2Fgit%40github.com%3AThomasObenaus%2Fsokar.git?ref=badge_shield)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=alert_status)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=code_smells)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=ncloc)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=security_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=sqale_index)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_sokar&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=ThomasObenaus_sokar)

## Overview

### Purpose

_Sokar is a generic alert based auto-scaler for cloud systems._

If you are running your microservices on a container orchestration system like [Nomad](https://www.nomadproject.io) or [kubernetes](https://kubernetes.io), then you are probably also need to scale them based on the load varying over time. The same situation applies if your system runs directly on AWS EC2 instances and thus they have to be scaled out with increased and scaled in with a reduction of the load. Usually the decision to scale is made based on metrics like current CPU/ RAM utilization or requests per second. But often you might want to use custom metrics like the length of a job-queue, the number of processed images per second or even a combination of those.

Here comes sokar into play. Sokar is a generic auto-scaler that makes scale up/ down decisions based on [scale alerts](doc/ScaleAlerts.md). He constantly evaluates the incoming scaling alerts, aggregates them and then scales the desired `ScaleObject` (i.e. microservice or an EC2 instance). Even if multiple metrics shall be taken into account for scaling the `ScaleObject`, those metrics just have to be expressed as scaling alerts and sokar will use them accordingly for scaling.

![doc/overview_coarse.png](doc/overview_coarse.png)

### Benefit

1. Possibility to combine multiple metrics to be taken into account for the scaling decisions. The impact of those metrics, expressed as [scale alerts](doc/ScaleAlerts.md), can be easily adjusted by configuring suitable weights.
2. Use the connectors to scale the actual `ScaleObject`. No need to implement the communication with the Container Orchestration System in this regard. The supported connectors can be found [here](doc/Connectors.md).
3. Configurable and ready to use capacity planning, providing separate cool downs for up and down scaling. Further more it is possible to select the [planning mode](doc/PlanningMode.md) which fits best for your workload.

### State

At the moment sokar is able to scale **Nomad jobs**, **Nomad instances (running on AWS)** and **AWS EC2 instances**.
For details about the changes see the [changelog](CHANGELOG.md).

## Build and Run

One can build and run sokar either in docker or directly on the host.

### On Host

```bash
# build
make build

# run (as scaler for nomad jobs)
make run.nomad-job
```

### Docker

```bash
# build
make docker.build

# pull it from docker hub
docker pull thobe/sokar

# run (as scaler for nomad jobs)
docker run -p 11000:11000 thobe/sokar:latest
# example
docker run -p 11000:11000 thobe/sokar:latest
```

For more configuration options and how to specify if sokar shall run as scaler for nomad jobs, nomad instances or AWS instances see [Config.md](config/Config.md).

## Features

- [Scheduled Scaling](doc/ScheduledScaling.md)
- [Dry Run Mode](doc/DryRunMode.md)
- [Several Planning Modes](doc/PlanningMode.md)
- [Several Connectors/ Scaling Targets](doc/Connectors.md)

## Links

- For configuration see [Config.md](config/Config.md)
- The API of sokar is specified at [openapi-spec.yaml](doc/openapi-spec.yaml)
- For metrics see [Metrics.md](Metrics.md)
- Structure and components of sokar are described at [Components.md](doc/Components.md)
- Example Grafana dashboards can be found at [Dashboards.md](dashboards/Dashboards.md)
