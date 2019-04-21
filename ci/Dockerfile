FROM golang:1.11.5-alpine AS golang

ARG PROJECT_PATH=github.com/thomasobenaus/sokar
ARG BINARY_NAME=sokar-bin

RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache make    

RUN mkdir -p /work

ENV GOPATH /work
ENV GOBIN /work/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

COPY . /work/src/${PROJECT_PATH}
WORKDIR /work/src/${PROJECT_PATH}

RUN make build



FROM alpine:latest

WORKDIR /app
COPY --from=golang /work/src/${PROJECT_PATH}/${BINARY_NAME} service


CMD [ "./sokar-bin" ]
