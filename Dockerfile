FROM alpine:3.6

RUN apk update && \
  apk add \
  ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-downstream-wait /bin/
ENTRYPOINT ["/bin/drone-downstream-wait"]
