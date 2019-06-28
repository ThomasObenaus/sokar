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

### Nomad

- This section contains the configuration parameters for nomad based scalers (i.e. job or data-center on AWS).

#### Mode

|         |                                                                                                                                                                                                                          |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| name    | mode                                                                                                                                                                                                                     |
| usage   | Scaling target mode is either job based or data-center (worker/ instance) based scaling. In data-center (dc) mode the nomad workers will be scaled. In job mode the number of allocations for this job will be adjusted. |
| type    | string (enum: job \| dc )                                                                                                                                                                                                |
| default | job                                                                                                                                                                                                                      |
| flag    | --sca.nomad.mode                                                                                                                                                                                                         |
| env     | SK_SCA_NOMAD_MODE                                                                                                                                                                                                        |

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

### [DEPRECATED] ScalingTarget

- Replaced by `--sca.nomad.mode`

|         |                                                                                                                                                                                                                          |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| name    | mode                                                                                                                                                                                                                     |
| usage   | Scaling target mode is either job based or data-center (worker/ instance) based scaling. In data-center (dc) mode the nomad workers will be scaled. In job mode the number of allocations for this job will be adjusted. |
| type    | string (enum: job \| dc )                                                                                                                                                                                                |
| default | job                                                                                                                                                                                                                      |
| flag    | --sca.mode                                                                                                                                                                                                               |
| env     | SK_SCA_MODE                                                                                                                                                                                                              |

## [DEPRECATED] Nomad

### [DEPRECATED] Server-Address

- Replaced by `--sca.nomad.server-address`

|         |                                            |
| ------- | ------------------------------------------ |
| name    | server-address                             |
| usage   | Specifies the address of the nomad server. |
| type    | string                                     |
| default | ""                                         |
| flag    | --nomad.server-address                     |
| env     | SK_NOMAD_SERVER_ADDRESS                    |

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
| usage   | Cycle/ frequency the ScaleAlertAggregator removes expired alerts.                                                                        |
| type    | List of value triplets (alert-name:alert-weight:alert-description). List elements are separated by a ';' and values are separated by '.' |
| default | ""                                                                                                                                       |
| example | --saa.scale-alerts="AlertA:1.0:An upscaling alert;AlertB:-1.5:A downscaling alert"                                                       |
| flag    | --saa.scale-alerts                                                                                                                       |
| env     | SK_SAA_SCALE_ALERTS                                                                                                                      |
