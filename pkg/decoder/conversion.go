package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
import "C"
import (
	"fmt"
	"unsafe"
//	"strings"
)

// octetStringToGoString takes in a ASN1 octet string and converts it to a Go string in hex
func octetStringToGoString(oString *C.OCTET_STRING_t) string {
	size := int(oString.size)
	str := ""
	for x := 0; x < size; x++ {
		octetByte := *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(oString.buf)) + uintptr(x)))
		str += fmt.Sprintf("%02X ", octetByte)
	}
	return str
}

// bitStringToGoString takes in a ASN1 bit string and converts it to a Go string in binary
func bitStringToGoString(bString *C.BIT_STRING_t) string {
	bitsUnused := uint64(bString.bits_unused)
	size := uint64(bString.size)
	body := uint8(*bString.buf)
	bits := int((size * 8) - bitsUnused)
	resStr := ""
	for x := 0; x < bits; x++ {
		if x := 0x80 & body; x == 128 {
			resStr += "1"
		} else {
			resStr += "0"
		}
		body = body << 1
	}
	return resStr
}