# Dry Run Mode

The following table shows how sokar behaves in case the dry run mode is activated.

| Feature                                                                                                                                             | Dry Run Mode Active                                                | Dry Run Mode Deactivated |
| :-------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------- | :----------------------- |
| Automatic Scaling                                                                                                                                   | Deactivated                                                        | Active                   |
| Manual Scaling                                                                                                                                      | Possible                                                           | Not Possible             |
| ScaleObjectWatcher                                                                                                                                  | Deactivated                                                        | Active                   |
| PlanedButSkippedScalingOpen<br>_(The metric `sokar_sca_planned_but_skipped_scaling_open`,<br>for more information see [Metrics.md](../Metrics.md))_ | Set to 1 if a scaling was skipped<br>Set to 0 after manual scaling | Stays 0                  |
