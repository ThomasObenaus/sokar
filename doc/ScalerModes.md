# Scaler Modes

Sokar supports to scale different scaling target types (nomad job, nomad node, AWS instances). For each of them sokar has to be configured to run in the according mode and depending on the scaling target type additional configuration parameters (e.g. credentials or a server address) have to be specified.

In this document the configuration of sokars scaler modes are described. As stated in [Config.md](../config/Config.md) the parameters can be set via environment variables, command line parameters or using a config file. For the sake of clarity a minimal config file, containing the parameters all tree modes have in common, and command line parameters are used here.

```yaml
scale-object:
  name: "fail-service"
  min: 1
  max: 10
cap:
  down-scale-cooldown: 20s
  up-scale-cooldown: 10s
  constant-mode:
    offset: 1
  scale-schedule:
    - days: "MON-FRI"
      start-time: 7:30
      end-time: 9:30
      min: 10
      max: 30
    - days: "SAT-SUN"
      start-time: 15:30
      end-time: 17:30
      min: 2
      max: 5
saa:
  no-alert-damping: 1.0
  up-thresh: 10.0
  down-thresh: -10.0
  eval-cycle: 1s
  eval-period-factor: 10
  cleanup-cycle: 60s
  alert-expiration-time: 5m
  scale-alerts:
    - name: "AlertA"
      weight: 1.5
      description: "Up alert"
    - name: "AlertB"
      weight: -1.5
      description: "Down alert"
logging:
  structured: false
  unix-ts: false
```

## Nomad Job

In this mode sokar is able to control the scale a nomad job by using the nomad api.

## Nomad Node (on AWS instances)

## AWS Instance
