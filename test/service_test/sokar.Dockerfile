FROM thobe/sokar:latest
LABEL maintainer="Thomas Obenaus"

WORKDIR /app
#COPY --from=golang /work/src/${PROJECT_PATH}/ci/config.yaml config.yaml

EXPOSE 11000

ENTRYPOINT [ "./sokar","--sca.nomad.mode=dc" ]