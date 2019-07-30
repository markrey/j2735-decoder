# J2735 Decoder

Go package for decoding J2735 CV2X standards. C source is compiled using ASN.1 compiler (https://github.com/vlm/asn1c). CGO is used for wrapping C calls from Go. Client interface uses MQTT to consume messages and decodes in either XML or JSON format for consumption. 

For build instructions, please see Dockerfile.
