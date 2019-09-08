FROM thobe/sokar:latest
LABEL maintainer="Thomas Obenaus"
ARG PROJECT_PATH=github.com/thomasobenaus/sokar

## Specifies the config.yaml that should be loaded and used for this test
ARG config_file

WORKDIR /app
COPY $config_file config.yaml

EXPOSE 11000

ENTRYPOINT [ "./sokar","--config-file=config.yaml"]