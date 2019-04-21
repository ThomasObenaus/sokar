
## The Build Image
FROM golang:1.11.5-alpine AS golang

ARG PROJECT_PATH=github.com/thomasobenaus/sokar
ARG BINARY_NAME=sokar

# Install needed tools
RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache make git

RUN mkdir -p /work

ENV GOPATH /work
ENV GOBIN /work/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

# Copy sources
COPY . /work/src/${PROJECT_PATH}
WORKDIR /work/src/${PROJECT_PATH}

# Build and test
RUN make build build_destination=/work name=${BINARY_NAME}

## The Run Image
FROM alpine:latest
LABEL maintainer="Thomas Obenaus"

ARG PROJECT_PATH=github.com/thomasobenaus/sokar
ARG BINARY_NAME=sokar

WORKDIR /app
COPY --from=golang /work/${BINARY_NAME} ${BINARY_NAME}
COPY --from=golang /work/src/${PROJECT_PATH}/ci/config.yaml config.yaml

ENTRYPOINT [ "./sokar","--config-file=config.yaml" ]
CMD [ "--help" ]
