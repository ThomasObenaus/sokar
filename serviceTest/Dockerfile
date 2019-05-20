
## The Build Image
FROM golang:1.11.5-alpine AS golang

ARG PROJECT_PATH=github.com/thomasobenaus/sokar

RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache gcc musl-dev

RUN mkdir -p /work

ENV GOPATH /work
ENV GOBIN /work/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

COPY . /work/src/${PROJECT_PATH}
WORKDIR /work/src/${PROJECT_PATH}/serviceTest

EXPOSE 4646
ENTRYPOINT [ "go","test","-sokar-address=http://localhost:11000" ]
