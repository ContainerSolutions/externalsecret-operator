FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN adduser -D externalconfig-operator
USER externalconfig-operator

ADD tmp/_output/bin/externalconfig-operator /usr/local/bin/externalconfig-operator
