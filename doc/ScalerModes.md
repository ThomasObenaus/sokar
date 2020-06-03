# Scaler Modes

Sokar supports to scale different scaling target types (nomad job, nomad node, AWS instances). For each of them sokar has to be configured to run in the according mode and depending on the scaling target type additional configuration parameters (e.g. credentials or a server address) have to be specified.

In this document the configuration of sokars scaler modes are described. As stated in [Config.md](../config/Config.md) the parameters can be set via environment variables, command line parameters or using a config file. For the sake of simplicity a minimal config file, containing the parameters all tree modes have in common, and command line parameters are used here.

## Minimal Configuration File

The configuration parameters have all useful default values to get sokar up and running quickly (at least in nomad-job mode using a locally running nomad server). Hence the minimal config we use here (`minimal.yaml`) just contains the definition of the alerts sokar should use to scale the according target.

```yaml
saa:
  scale-alerts:
    - name: "AlertA"
      weight: 1.5
      description: "Up alert"
    - name: "AlertB"
      weight: -1.5
      description: "Down alert"
```

## Nomad Job

In this mode sokar is able to control the scale a nomad job by using the nomad api. As `scale-object` the [thobe/fail_service](https://hub.docker.com/r/thobe/fail_service) docker image is used. Hence it has to be deployed to nomad already (therefore [multi-group.nomad.nomad](../examples/multi-group.nomad) can be used).

To run sokar in scaler mode `nomad-job` one just has to start sokar providing the minimal configuration file and the name of the `scale-object`, which is `fail-service`.

```bash
# start sokar to scale the nomad job named fail-service
./sokar-bin --config-file=examples/config/minimal.yaml --scale-object.name="fail-service"
```

## Nomad Node (on AWS instances)

## AWS Instance
