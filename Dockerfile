FROM golang:1.14-alpine AS build-stage

MAINTAINER contact@hexhibit.xyz

RUN apk add --update --no-cache build-base\
    && apk add git openssh

RUN mkdir /agor
WORKDIR /agor/

COPY go.mod .

RUN go mod download

COPY ./ ./


RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -v -o agor /agor/


# production stage
FROM alpine:3.5
MAINTAINER christian.schlatter@ionesoft.ch

RUN apk --update add bash

#Copy binary
COPY --from=build-stage /agor/agor /etc/agor/
COPY ./startup.sh /etc/agor/
COPY ./assets/ /etc/agor/assets/
COPY ./config.yaml /etc/agor/config/default.yaml

RUN chmod +x /etc/agor/startup.sh

ENTRYPOINT ["/bin/bash", "-c", "/etc/agor/startup.sh"]