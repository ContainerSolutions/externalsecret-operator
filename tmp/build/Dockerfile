FROM alpine:3.6

RUN adduser -D externalconfig-operator
USER externalconfig-operator

ADD tmp/_output/bin/externalconfig-operator /usr/local/bin/externalconfig-operator
