FROM golang:1.12-alpine as builder
RUN apk update && \
    apk add --update make git
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /src
RUN go build -o test-agent ./cmd/mqtt-pub/*.go

FROM golang:1.12-alpine
RUN apk update
WORKDIR /app
COPY --from=builder /src/pkg/decoder/samples/logs/bsm-10-23.log .
COPY --from=builder /src/pkg/decoder/samples/spat.uper .
COPY --from=builder /src/test-agent .
ENTRYPOINT [ "/app/test-agent" ]