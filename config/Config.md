# Configuration

Sokar can be configured either through a config-file, environment variables or command-line parameters.
The order they are applied is:

1. Default values are overwritten by
2. Parameters defined in the config-file ([full example](examples/config/full.yaml)), which are overwritten by
3. Environment variables, which are overwritten by
4. Command-Line parameters

## Global

### Config-File

|         |                                                                      |
| ------- | -------------------------------------------------------------------- |
| name    | config-file                                                          |
| usage   | Specifies the full path and name of the configuration file for sokar |
| type    | string                                                               |
| default | ""                                                                   |
| flag    | --config-file                                                        |
| env     | -                                                                    |

### Dry-Run

|         |                                                                                                                                      |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| name    | dry-run                                                                                                                              |
| usage   | If true, then sokar won't execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed. |
| type    | bool                                                                                                                                 |
| default | false                                                                                                                                |
| flag    | --dry-run                                                                                                                            |
| env     | SK_DRY_RUN                                                                                                                           |

### Port

|         |                                |
| ------- | ------------------------------ |
| name    | port                           |
| usage   | Port where sokar is listening. |
| type    | uint                           |
| default | 11000                          |
| flag    | --port                         |
| env     | SK_PORT                        |

## Scaler

### ScalingObjectWatcherInterval

|         |                                                                                                |
| ------- | ---------------------------------------------------------------------------------------------- |
| name    | watcher-interval                                                                               |
| usage   | The interval the Scaler will check if the scalingObject count still matches the desired state. |
| type    | duration                                                                                       |
| default | 5s                                                                                             |
| flag    | --sca.watcher-interval                                                                         |
| env     | SK_SCA_WATCHER_INTERVAL                                                                        |

### Mode

|         |                                                                                                                                                                                                                                                                                                                                                        |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| name    | mode                                                                                                                                                                                                                                                                                                                                                   |
| usage   | Scaling target mode is either nomad-job based, aws EC2 or nomad data-center (worker/ instance) based scaling. In data-center (`noamd-dc`) mode the nomad workers will be scaled. In `nomad-job` mode the number of allocations for this job will be adjusted. In `aws-ec2` mode AWS instances will be scaled adjusting the according AutoScalingGroup. |
| type    | string (enum: nomad-job \| nomad-dc \| aws-ec2 )                                                                                                                                                                                                                                                                                                       |
| default | job                                                                                                                                                                                                                                                                                                                                                    |
| flag    | --sca.mode                                                                                                                                                                                                                                                                                                                                             |
| env     | SK_SCA_MODE                                                                                                                                                                                                                                                                                                                                            |

### AWS EC2

- This section contains the configuration parameters for AWS EC2 based scaling.

#### Profile

|         |                                                                                                                                                                                                                                                                                                                                                    |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name    | profile                                                                                                                                                                                                                                                                                                                                            |
| usage   | This parameter represents the name of the aws profile that shall be used to access the resources to scale the data-center. This parameter is optional. If it is empty the instance where sokar runs on has to have enough permissions to access the resources (ASG) for scaling. In this case the AWSRegion parameter has to be specified as well. |
| type    | string                                                                                                                                                                                                                                                                                                                                             |
| default | ""                                                                                                                                                                                                                                                                                                                                                 |
| flag    | --sca.aws-ec2.profile                                                                                                                                                                                                                                                                                                                              |
| env     | SK_SCA_AWS_EC2_PROFILE                                                                                                                                                                                                                                                                                                                             |

#### Region

|         |                                                                                                                                                                                    |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name    | region                                                                                                                                                                             |
| usage   | This is an optional parameter and is regarded only if the parameter AWSProfile is empty. The AWSRegion has to specify the region in which the data-center to be scaled resides in. |
| type    | string                                                                                                                                                                             |
| default | ""                                                                                                                                                                                 |
| flag    | --sca.aws-ec2.region                                                                                                                                                               |
| env     | SK_SCA_AWS_EC2_REGION                                                                                                                                                              |

#### ASGTagKey

|         |                                                                                                                                  |
| ------- | -------------------------------------------------------------------------------------------------------------------------------- |
| name    | asg_tag_key                                                                                                                      |
| usage   | This parameter specifies which tag on an AWS AutoScalingGroup shall be used to find the ASG that should be automatically scaled. |
| type    | string                                                                                                                           |
| default | "scale-object"                                                                                                                   |
| flag    | --sca.aws-ec2.asg-tag-key                                                                                                        |
| env     | SK_SCA_AWS_EC2_ASG_TAG_KEY                                                                                                       |

### Nomad

- This section contains the configuration parameters for nomad based scalers (i.e. job or data-center on AWS).

#### Server-Address

|         |                                            |
| ------- | ------------------------------------------ |
| name    | server-address                             |
| usage   | Specifies the address of the nomad server. |
| type    | string                                     |
| default | ""                                         |
| flag    | --sca.nomad.server-address                 |
| env     | SK_SCA_NOMAD_SERVER_ADDRESS                |

#### Data-Center AWS

- The parameters in this section are used to configure the scaler that is used to scale a data-center hosted on AWS

##### Profile

|         |                                                                                                                                                                                                                                                                                                                                                    |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name    | profile                                                                                                                                                                                                                                                                                                                                            |
| usage   | This parameter represents the name of the aws profile that shall be used to access the resources to scale the data-center. This parameter is optional. If it is empty the instance where sokar runs on has to have enough permissions to access the resources (ASG) for scaling. In this case the AWSRegion parameter has to be specified as well. |
| type    | string                                                                                                                                                                                                                                                                                                                                             |
| default | ""                                                                                                                                                                                                                                                                                                                                                 |
| flag    | --sca.nomad.dc-aws.profile                                                                                                                                                                                                                                                                                                                         |
| env     | SK_SCA_NOMAD_DC_AWS_PROFILE                                                                                                                                                                                                                                                                                                                        |

##### Region

|         |                                                                                                                                                                                    |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name    | region                                                                                                                                                                             |
| usage   | This is an optional parameter and is regarded only if the parameter AWSProfile is empty. The AWSRegion has to specify the region in which the data-center to be scaled resides in. |
| type    | string                                                                                                                                                                             |
| default | ""                                                                                                                                                                                 |
| flag    | --sca.nomad.dc-aws.region                                                                                                                                                          |
| env     | SK_SCA_NOMAD_DC_AWS_REGION                                                                                                                                                         |

##### Timeout for Instance Termination

|         |                                                                                                                                              |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| name    | instance-termination-timeout                                                                                                                 |
| usage   | The maximum time the instance termination will be monitored before assuming that this action (instance termination due to downscale) failed. |
| type    | duration                                                                                                                                     |
| default | 10m                                                                                                                                          |
| flag    | --sca.nomad.dc-aws.instance-termination-timeout                                                                                              |
| env     | SK_SCA_NOMAD_DC_AWS_INSTANCE_TERMINATION_TIMEOUT                                                                                             |

## ScaleObject

### Name

|         |                                      |
| ------- | ------------------------------------ |
| name    | name                                 |
| usage   | The name of the object to be scaled. |
| type    | string                               |
| default | ""                                   |
| flag    | --scale-object.name                  |
| env     | SK_SCALE_OBJECT_NAME                 |

### Min

|         |                                               |
| ------- | --------------------------------------------- |
| name    | min                                           |
| usage   | The minimum count of the object to be scaled. |
| type    | uint                                          |
| default | 1                                             |
| flag    | --scale-object.min                            |
| env     | SK_SCALE_OBJECT_MIN                           |

### Max

|         |                                               |
| ------- | --------------------------------------------- |
| name    | max                                           |
| usage   | The maximum count of the object to be scaled. |
| type    | uint                                          |
| default | 10                                            |
| flag    | --scale-object.max                            |
| env     | SK_SCALE_OBJECT_MAX                           |

## CapacityPlanner

### Scale Schedule

|          |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name     | scale-schedule                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| usage    | Specifies time ranges within which it is ensured that the ScaleObject is scaled to at least min and not more than max. The `min`/ `max` values specified in this schedule have lower priority than the `--scale-object.min`/ `--scale-object.max`. This means the sokar will ensure that the `--scale-object.min`/ `--scale-object.max` are not violated no matter what is specified in the schedule.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| type     | The schedule is a list of scale schedule entries.</br>These entries are separated by `|`:</br><ul><li>`"<entry>|<entry>|..."`</li></ul>One entry consist of five parts `<days> <start-time> <end-time> <scale-range>`. The parts are separated by `<space>`: <ol><li>**`<days>`**<ul><li>Specifies at which weekdays the schedule shall be active.</li><li>Valid values: `MON,TUE,WED,THU,FRI,SAT,SUN` and `0,1,2,3,4,5,6`</li><li>Ranges are allowed e.g. `MON-FRI`, but have to be ascending (e.g. `FRI-MON` is invalid)</li><li>The wildcard `*` is also allowed and means any day. But `*-*` or `<day>-*` (e.g. `FRI-*`) are not valid.</li></ul></li><li>**`<start-time>`**<ul><li>Specifies the time at which the schedule begins.<li>Format: `<hour>:<minute>`. Where hour is a number between 0 and 23 and minute a number between 0 and 59.</li><li>The minute qualifier is optional. Thus instead of specifying `13:00`, `13` is sufficient.</li><li>The wildcard `*` is also allowed and means start of the day (0:00 - midnight).</li></ul></li><li>**`<end-time>`**<ul><li>Specifies the time at which the schedule ends.<li>Format: `<hour>:<minute>`. Where hour is a number between 0 and 23 and minute a number between 0 and 59.</li><li>The minute qualifier is optional. Thus instead of specifying `13:00`, `13` is sufficient.</li><li>The wildcard `*` is also allowed and means end of the day (0:00 - midnight).</li></ul></li><li>**`<scale-range>`**<ul><li>Specifies the range within which the scale of the scale object shall be kept during this schedule.<li>Format: `<min-scale>-<max-scale>`, where both scale values are `uint`.</li><li>The wildcard `*` is also allowed and means unbound. For example `*-10` means the min scale is not bound whereas the max is set to 10.</li><li>It is allowed to just specify `*` instead of `*-*` if both, min- and max-scale shall be unbound. Even though it makes no sense, since no scheduled scaling would be done in this case.</li></ul></li></ol>**Rules:** <ol><li>In case no scheduled scale entry is active according to the specified time ranges, no scheduled scaling will be done.</li><li>The schedule entries must not overlap.</li></ol> |
| examples | Rush hour peaks: <ul><li>`--cap.scale-schedule=MON-FRI 7 9 10-30|MON-FRI 14 18 10-30`</li><li>Scale between 10 and 30 during 7am-9am and 4pm-6pm at businessdays.</li><li>At the rest of the day and on weekend no scheduled scaling will be applied.</li></ul>Over the weekend:<ul><li>`--cap.scale-schedule=SAT-SUN * * 2-5`</li><li>Scale between 2 and 5 during weekend.</li></ul>Half hours:<ul><li>`--cap.scale-schedule=* 7:30 9:30 10-30`</li><li>Scale between 10 and 30 during 7:30am-9:30am.</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| default  | "" (no schedule)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| flag     | --cap.scale-schedule                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| env      | SK_CAP_SCALE_SCHEDULE                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

### Down Scale Cooldown

|         |                                                          |
| ------- | -------------------------------------------------------- |
| name    | down-scale-cooldown                                      |
| usage   | The time sokar waits between downscaling actions at min. |
| type    | duration                                                 |
| default | 20s                                                      |
| flag    | --cap.down-scale-cooldown                                |
| env     | SK_CAP_DOWN_SCALE_COOLDOWN                               |

### Up Scale Cooldown

|         |                                                        |
| ------- | ------------------------------------------------------ |
| name    | up-scale-cooldown                                      |
| usage   | The time sokar waits between upscaling actions at min. |
| type    | duration                                               |
| default | 20s                                                    |
| flag    | --cap.up-scale-cooldown                                |
| env     | SK_CAP_UP_SCALE_COOLDOWN                               |

### Constant Mode

#### Enable

|         |                                                                                                                  |
| ------- | ---------------------------------------------------------------------------------------------------------------- |
| name    | enable                                                                                                           |
| usage   | Enable/ disable the constant mode of the CapacityPlanner. Only one of the modes can be enabled at the same time. |
| type    | bool                                                                                                             |
| default | true                                                                                                             |
| flag    | --cap.constant-mode.enable                                                                                       |
| env     | SK_CAP_CONSTANT_MODE_ENABLE                                                                                      |

#### Offset

|         |                                                                                                                                 |
| ------- | ------------------------------------------------------------------------------------------------------------------------------- |
| name    | offset                                                                                                                          |
| usage   | The constant offset value that should be used to increment/ decrement the count of the scale-object. Only values > 0 are valid. |
| type    | uint                                                                                                                            |
| default | 1                                                                                                                               |
| flag    | --cap.constant-mode.offset                                                                                                      |
| env     | SK_CAP_CONSTANT_MODE_OFFSET                                                                                                     |

### Linear Mode

#### Enable

|         |                                                                                                                |
| ------- | -------------------------------------------------------------------------------------------------------------- |
| name    | enable                                                                                                         |
| usage   | Enable/ disable the linear mode of the CapacityPlanner. Only one of the modes can be enabled at the same time. |
| type    | bool                                                                                                           |
| default | false                                                                                                          |
| flag    | --cap.linear-mode.enable                                                                                       |
| env     | SK_CAP_LINEAR_MODE_ENABLE                                                                                      |

#### ScaleFactorWeight

|         |                                                                                                      |
| ------- | ---------------------------------------------------------------------------------------------------- |
| name    | scale-factor-weight                                                                                  |
| usage   | This weight is used to adjust the impact of the scaleFactor during capacity planning in linear mode. |
| type    | float                                                                                                |
| default | 0.5                                                                                                  |
| flag    | --cap.linear-mode.scale-factor-weight                                                                |
| env     | SK_CAP_LINEAR_MODE_SCALE_FACTOR_WEIGHT                                                               |

## Logging

### Structured

|         |                                |
| ------- | ------------------------------ |
| name    | structured                     |
| usage   | Use structured logging or not. |
| type    | bool                           |
| default | false                          |
| flag    | --logging.structured           |
| env     | SK_LOGGING_STRUCTURED          |

### Unix Timestamp

|         |                                                    |
| ------- | -------------------------------------------------- |
| name    | unix-ts                                            |
| usage   | Use Unix-Timestamp representation for log entries. |
| type    | bool                                               |
| default | false                                              |
| flag    | --logging.unix-ts                                  |
| env     | SK_LOGGING_UNIX_TS                                 |

### NoColor

|         |                                                 |
| ------- | ----------------------------------------------- |
| name    | no-color                                        |
| usage   | If true colors in log out-put will be disabled. |
| type    | bool                                            |
| default | false                                           |
| flag    | --logging.no-color                              |
| env     | SK_LOGGING_NO_COLOR                             |

### Level

|         |                                                                                                    |
| ------- | -------------------------------------------------------------------------------------------------- |
| name    | level                                                                                              |
| usage   | The level that should be used for logs. Valid entries are debug, info, warn, error, fatal and off. |
| type    | string                                                                                             |
| default | info                                                                                               |
| flag    | --logging.level                                                                                    |
| env     | SK_LOGGING_LEVEL                                                                                   |

## ScaleAlertAggregator

### Alert Expiration Time

|         |                                                                                                                                                              |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| name    | alert-expiration-time                                                                                                                                        |
| usage   | Defines after which time an alert will be pruned if he did not get updated again by the ScaleAlertEmitter, assuming that the alert is not relevant any more. |
| type    | Duration                                                                                                                                                     |
| default | 10m                                                                                                                                                          |
| flag    | --saa.alert-expiration-time                                                                                                                                  |
| env     | SK_SAA_ALERT_EXPIRATION_TIME                                                                                                                                 |

### No Alert Damping

|         |                                                                                          |
| ------- | ---------------------------------------------------------------------------------------- |
| name    | no-alert-damping                                                                         |
| usage   | Damping used in case there are currently no alerts firing (neither down- nor upscaling). |
| type    | float                                                                                    |
| default | 1.0                                                                                      |
| flag    | --saa.no-alert-damping                                                                   |
| env     | SK_SAA_NO_ALERT_DAMPING                                                                  |

### Up Scale Threshold

|         |                                  |
| ------- | -------------------------------- |
| name    | up-thresh                        |
| usage   | Threshold for a upscaling event. |
| type    | float                            |
| default | 10.0                             |
| flag    | --saa.up-thresh                  |
| env     | SK_SAA_UP_THRESH                 |

### Down Scale Threshold

|         |                                    |
| ------- | ---------------------------------- |
| name    | down-thresh                        |
| usage   | Threshold for a downscaling event. |
| type    | float                              |
| default | 10.0                               |
| flag    | --saa.down-thresh                  |
| env     | SK_SAA_DOWN_THRESH                 |

### Evaluation Cycle

|         |                                                                                                 |
| ------- | ----------------------------------------------------------------------------------------------- |
| name    | eval-cycle                                                                                      |
| usage   | Cycle/ frequency the ScaleAlertAggregator evaluates the weights of the currently firing alerts. |
| type    | Duration                                                                                        |
| default | 1s                                                                                              |
| flag    | --saa.eval-cycle                                                                                |
| env     | SK_SAA_EVAL_CYCLE                                                                               |

### Evaluation Period Factor

|         |                                                                                                                                   |
| ------- | --------------------------------------------------------------------------------------------------------------------------------- |
| name    | eval-period-factor                                                                                                                |
| usage   | EvaluationPeriodFactor is used to calculate the evaluation period (evaluationPeriod = evaluationCycle \* evaluationPeriodFactor). |
| type    | uint                                                                                                                              |
| default | 10                                                                                                                                |
| flag    | --saa.eval-period-factor                                                                                                          |
| env     | SK_SAA_EVAL_PERIOD_FACTOR                                                                                                         |

### Cleanup Cycle

|         |                                                                   |
| ------- | ----------------------------------------------------------------- |
| name    | cleanup-cycle                                                     |
| usage   | Cycle/ frequency the ScaleAlertAggregator removes expired alerts. |
| type    | Duration                                                          |
| default | 60s                                                               |
| flag    | --saa.cleanup-cycle                                               |
| env     | SK_SAA_CLEANUP_CYCLE                                              |

### Scale Alerts

|         |                                                                                                                                          |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| name    | scale-alerts                                                                                                                             |
| usage   | The alerts that should be used for scaling (up/down) the scale-object.                                                                   |
| type    | List of value triplets (alert-name:alert-weight:alert-description). List elements are separated by a ';' and values are separated by '.' |
| default | ""                                                                                                                                       |
| example | --saa.scale-alerts="AlertA:1.0:An upscaling alert;AlertB:-1.5:A downscaling alert"                                                       |
| flag    | --saa.scale-alerts                                                                                                                       |
| env     | SK_SAA_SCALE_ALERTS                                                                                                                      |
