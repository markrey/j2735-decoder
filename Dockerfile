FROM golang:1.12-alpine as builder
RUN apk update && \
    apk add --update make git musl-dev gcc
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /src/pkg/decoder/c
RUN make -f converter-example.mk clean
RUN make -f converter-example.mk libasncodec.a
WORKDIR /src
RUN go build -o sdmap-agent ./cmd/sdmap-agent/*.go

FROM golang:1.12-alpine
RUN apk update
WORKDIR /app
COPY --from=builder /src/sdmap-agent .
ENTRYPOINT [ "/app/sdmap-agent" ]
