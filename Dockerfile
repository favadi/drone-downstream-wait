# build stage
FROM golang:alpine AS build-env

ENV GOPATH=/go
ADD . /go/src/github.com/favadi/drone-downstream-wait
WORKDIR /go
RUN go install -v github.com/favadi/drone-downstream-wait

# final stage
FROM alpine:3.7
RUN apk update && \
  apk add \
  ca-certificates && \
  rm -rf /var/cache/apk/*

COPY --from=build-env /go/bin/drone-downstream-wait /bin/
ENTRYPOINT ["/bin/drone-downstream-wait"]
