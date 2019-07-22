FROM golang:1.12-alpine
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/decoder-client/