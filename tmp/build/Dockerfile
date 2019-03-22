FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN adduser -D externalsecret-operator
USER externalsecret-operator

ADD tmp/_output/bin/externalsecret-operator /usr/local/bin/externalsecret-operator
