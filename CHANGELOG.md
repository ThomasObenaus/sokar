# Changelog

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
