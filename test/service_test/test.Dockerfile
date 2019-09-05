## The Build Image
FROM golang:1.11.5-alpine AS golang

ARG PROJECT_PATH=github.com/thomasobenaus/sokar/test/service_test

# Install needed tools
RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache make git gcc musl-dev

RUN mkdir -p /work
RUN mkdir -p /work/bin

ENV GOPATH /work
ENV GOBIN /work/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

# Copy sources
COPY . /work/src/${PROJECT_PATH}
WORKDIR /work/src/${PROJECT_PATH}

RUN make deps-install

# Just an empty cmd here. This docker file is intended to be called with 'make <make-target>'.
# Where the make target is the test that should be executed.
CMD [""]