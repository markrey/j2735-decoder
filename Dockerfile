FROM golang
ADD . /go/src/uper_decoder
RUN go get /go/src/uper_decoder/cmd/decoder-client
RUN go install /go/src/uper_decoder/cmd/decoder-client
ENTRYPOINT [ "/go/bin/decoder-client" ]
