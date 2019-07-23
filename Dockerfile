FROM golang:1.12-alpine as builder
RUN apk update && \
    apk add --update make git musl-dev gcc
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /src/pkg/decoder/c
RUN make -f converter-example.mk libasncodec.a
WORKDIR /src
RUN go build -o j2735-decoder ./cmd/decoder-client/*.go

FROM golang:1.12-alpine
RUN apk update
WORKDIR /app
COPY --from=builder /src/j2735-decoder .
ENTRYPOINT [ "/app/j2735-decoder" ]