FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /go/src/github.com/caioreix

COPY . .
RUN apk add --no-cache git tree make
RUN cd cmd/api && make build
RUN mv cmd/api/api-build build/api

FROM alpine
WORKDIR /root
COPY --from=builder /go/src/github.com/caioreix/build/api/api-build .
CMD ["./api-build"]
