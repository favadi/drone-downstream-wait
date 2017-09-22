FROM alpine:3.6

RUN apk update && \
  apk add \
  ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-deploy /bin/
ENTRYPOINT ["/bin/drone-deploy"]
