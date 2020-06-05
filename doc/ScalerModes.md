# Scaler Modes

In this document the configuration of sokars' scaler modes are described.

Each mode differs in the type of the scaling target. Sokar supports to scale different scaling target types (nomad job, nomad node, AWS EC2 instances). For each of them sokar has to be configured to run in the according mode and depending on the scaling target type, additional configuration parameters (e.g. credentials or a server address) have to be specified.

For the sake of simplicity a minimal config file, containing the parameters all tree modes have in common, and command line parameters are used here. But as stated in [Config.md](../config/Config.md) the parameters can be set via environment variables, command line parameters or using a config file as well.

## Prerequisites

### Minimal Configuration File

The configuration parameters have all useful default values to get sokar up and running quickly (at least in nomad-job mode using a locally running nomad server). Hence the minimal config used here (`minimal.yaml`) contains only the definition of the alerts sokar should use to scale the according scale object. The definition contains one alert for scaling up (_AlertA_) and one for scaling down (_AlertB_).

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

### Issue the Alert to Trigger the Scaling

For testing the according scaler mode one has to issue a scaling alert to sokar. This can be done with the following curl calls.
As defined in the minimal config `AlertA` would lead to an up- and `AlertB` to a down-scaling of the scale object.

```bash
# Issue a request to signal that 'AlertA' is active (firing) ==> UP-SCALING
curl --request POST 'http://localhost:11000/api/alerts' \
--header 'Content-Type: application/json' \
--data-raw '{
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "AlertA"
      }
    }
  ]
}'

# Issue a request to signal that 'AlertB' is active (firing) ==> DOWN-SCALING
curl --request POST 'http://localhost:11000/api/alerts' \
--header 'Content-Type: application/json' \
--data-raw '{
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "AlertB"
      }
    }
  ]
}'
```

## Nomad Job

In this mode sokar is able to control the scale a nomad job by using the nomad api. As `scale-object` the [thobe/fail_service](https://hub.docker.com/r/thobe/fail_service) docker image is used. Hence it has to be deployed to nomad already (therefore the nomad job file [multi-group.nomad.nomad](../examples/multi-group.nomad) can be used).

To run sokar in scaler mode `nomad-job` one just has to start sokar providing the minimal configuration file and the name of the `scale-object`, which is `fail-service`.

```bash
# start sokar to scale the nomad job named fail-service
./sokar-bin --config-file=examples/config/minimal.yaml --scale-object.name="fail-service"
```

## Nomad Data-Center (on AWS instances)

In this mode sokar is able to control the scale of the nomad nodes but only in case they are running on AWS EC2 instances that are managed by an AWS AutoScalingGroup (ASG). To be precise sokar manages the count of nodes that are running for one nomad data-center.

This means if due to the currently active scaling alerts an **up-scaling** is necessary, sokar will create new AWS EC2 instances by incrementing the ASG by the calculated amount of additionally needed instances.

In case of a needed **down-scaling** sokar will do the following steps as long as the calculated amount for down-scaling is satisfied:

1. Select the nomad node which would be the best candidate for termination. That is the instance with least capacity utilization (memory, cpu and then disk).
2. Drain the nomad node, which forces nomad to migrate jobs running on that node to other nodes.
3. Terminate the EC2 instance that is hosting the node that was selected for termination.

To run sokar in scaler mode `nomad-dc` one just has to start sokar providing the minimal configuration file and the following parameters:

- The scaler-mode: `--sca.mode=nomad-dc`.
- The name of the scale-object which is, in this mode, the name of the data-center: `--scale-object.name=my-data-center`.
- The address of nomad running on AWS instances: `--sca.nomad.server-address="http://nomad.example.com`.
- The [AWS profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html) that has the needed permissions to modify the AutoScalingGroup nomads data-center is running on: `--sca.nomad.dc-aws.profile=my-profile`.
- The AWS region where the nomad data-center is running: `--sca.nomad.dc-aws.region=eu-central-1`.

```bash
# start sokar to scale the nomad data-center named 'my-data-center'
./sokar-bin --config-file=examples/config/minimal.yaml \
  --sca.mode=nomad-dc \
  --scale-object.name=my-data-center \
  --sca.nomad.server-address="http://nomad.example.com" \
  --sca.nomad.dc-aws.profile=my-profile \
  --sca.nomad.dc-aws.region=eu-central-1
```

**Important Hint:** To enable sokar to scale the data-center it is necessary to tag the AWS AutoScalingGroup with the name of the nomad data-center.
The tag on that ASG has to have the key `scale-object` and the value has to be the name of the nomad data-center (e.g. `my-data-center`). Otherwise sokar is not able to identify the ASG that manages the instances the data-center is running on.

## AWS Instance

In this mode sokar is able to control the scale of AWS EC2 instances that are managed by an AWS AutoScalingGroup (ASG). To be precise sokar manages the count of instances based on the scale-alerts that are issued.

To run sokar in scaler mode `aws-ec2` one just has to start sokar providing the minimal configuration file and the following parameters:

- The scaler-mode: `--sca.mode=aws-ec2`.
- The name of the scale-object which is in this mode the name of the data-center: `--scale-object.name=my-autoscaling-group`.
- The [AWS profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html) that has needed permissions to modify the AutoScalingGroup: `--sca.aws-ec2.profile=my-profile`.
- The AWS region where the EC2 instances are running: `--sca.aws-ec2.region=eu-central-1`.

```bash
# start sokar to scale the nomad data-center named 'my-data-center'
./sokar-bin --config-file=examples/config/minimal.yaml \
  --sca.mode=aws-ec2 \
  --scale-object.name=my-autoscaling-group \
  --sca.aws-ec2.profile=my-profile \
  --sca.aws-ec2.region=eu-central-1
```

**Important Hint:** To enable sokar to scale the AWS AutoScalingGroup it is necessary to tag the AWS AutoScalingGroup.
The tag on that ASG has to have the key `scale-object` and the value has to be the name that is referred with the parameter `--scale-object.name` (e.g. `my-autoscaling-group`). Otherwise sokar is not able to identify the ASG that manages the EC2 instances.
