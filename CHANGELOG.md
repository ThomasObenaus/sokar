# Changelog

## v0.0.6 (2019-05-27)

- Data-Center Scaling: With [#65](https://github.com/ThomasObenaus/sokar/issues/60) and [#60](https://github.com/ThomasObenaus/sokar/issues/65) the functionality to scale not only jobs but also a datacenter was added. In the currently implementation an up-scaling and a **hard down-scaling is implemented**. This means **during down-scaling NO node draining will be applied**, the datacenter instances are directly terminated.
- Config: With [#64](https://github.com/ThomasObenaus/sokar/issues/64) the following config-parameters where replaced by new ones and are thus deprecated from now on:
  - `--dummy-scaling-target` is replaced by `--sca.mode=dc`. With `--sca.mode` the scaler mode can be configured to scale either a job (`--sca.mode=job`) or a datacenter (`--sca.mode=dc`).
  - `--job.name` is replaced by `--scale-object.name`
  - `--job.min` is replaced by `--scale-object.min`
  - `--job.max` is replaced by `--scale-object.max`
- Docu: With [#54](https://github.com/ThomasObenaus/sokar/issues/54) metrics documentation was added to [Metrics.md](https://github.com/ThomasObenaus/sokar/blob/master/Metrics.md).
- Fix: With [#57](https://github.com/ThomasObenaus/sokar/issues/57) an issue was fixed that sokar did ignored the down-scaling threshold.

## v0.0.5 (2019-05-12)

- Config: The expiration time used to prune scale alerts that are not updated for a long time is now configurable (--saa.alert-expiration-time, default: 10m).
- API: The scale-by end-points are enabled for dry-run mode only. This means in case sokar runs not in dry-run mode these and-points will respond an error and won't do any action.
- Robustness: Sokar now immediately adjusts the scale count at startup in case the deployed scale violates the min/ max bounds of the job.
- Robustness: Sokar now immediately adjusts the scale count to its expected count in case another source has adjusted the job count in between.

## v0.0.4 (2019-04-22)

- CI/CD: Sokar is now dockerized and available on [docker hub](https://hub.docker.com/r/thobe/sokar). It can be pulled via `docker pull thobe/sokar`.
- Config: A http end-point was added that provides the configuration used by sokar. It can be requested via GET request at `/api/config`.
- General: A http end-point was added that provides the build information of sokar. It can be requested via GET request at `/api/build`.

  - Example response:

    ```json
    {
      "version": "v0.0.3-17-ga9e9d18",
      "build_time": "2019-04-22_16-36-49",
      "revision": "a9e9d18_dirty",
      "branch": "f/34_configuration_endpoint"
    }
    ```

## v0.0.3 (2019-04-14)

- Cli/ Config: All parameters can now be configured/ set via cli, environment-variable or config-file. [#38](https://github.com/ThomasObenaus/sokar/issues/38)
- Metrics: Added basic metrics to all components. [#26](https://github.com/ThomasObenaus/sokar/issues/26)
- Dry-Run: A dry run mode was added. This can be used to test sokars decisions. [#31](https://github.com/ThomasObenaus/sokar/issues/31)
- Scale-By: A scale-by end-point was added. It can be used to manually scale (relative based on percentage or value). [#27](https://github.com/ThomasObenaus/sokar/issues/27)
