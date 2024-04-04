FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /go/src/github.com/caioreix

COPY . .
RUN apk add --no-cache git tree make
RUN cd cmd/worker && make build
RUN mv cmd/worker/worker-build build/worker

FROM alpine
WORKDIR /root
COPY --from=builder /go/src/github.com/caioreix/build/worker/worker-build .
CMD ["./worker-build"]
