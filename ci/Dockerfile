
## The Build Image
FROM golang:1.11.5-alpine AS golang

ARG PROJECT_PATH=github.com/thomasobenaus/sokar
ARG BINARY_NAME=sokar
ARG SRC_DIR=.

# Install needed tools
RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache make git gcc musl-dev

RUN mkdir -p /work
RUN mkdir -p /work/bin

ENV GOPATH /work
ENV GOBIN /work/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

## Copy sources
COPY ${SRC_DIR} /work/src/${PROJECT_PATH}
WORKDIR /work/src/${PROJECT_PATH}

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

# Deps
RUN make deps-install

# Test
RUN make tools test

# Build
RUN make build build_destination=/work name=${BINARY_NAME}

## The Run Image
FROM alpine:3.9
LABEL maintainer="Thomas Obenaus"

ARG PROJECT_PATH=github.com/thomasobenaus/sokar
ARG BINARY_NAME=sokar

WORKDIR /app
COPY --from=golang /work/${BINARY_NAME} ${BINARY_NAME}
COPY --from=golang /work/src/${PROJECT_PATH}/ci/config.yaml config.yaml

EXPOSE 11000

ENTRYPOINT [ "./sokar","--config-file=config.yaml" ]
CMD [ "--help" ]