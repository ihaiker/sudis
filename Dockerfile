FROM golang:1.13.6-alpine3.11 as builder

ENV GO111MODULE="on"
ARG LDFLAGS=""

ADD . /sudis
WORKDIR /sudis

RUN go get -u github.com/shuLhan/go-bindata/cmd/go-bindata
RUN go generate generator.go
RUN apk add --no-cache make git build-base
RUN go build -tags bindata -ldflags "${LDFLAGS}" -o sudis sudis.go


FROM alpine:3.11
MAINTAINER Haiker ni@renzhen.la

WORKDIR /opt/sudis

ADD ./libs/config/sudis.example.yaml /etc/sudis/sudis.yaml
COPY --from=builder /sudis/sudis /opt/sudis/sudis

ENV PATH="$PATH:/opt/sudis"
EXPOSE 5983 5984
VOLUME /etc/sudis

ENTRYPOINT ["sudis"]