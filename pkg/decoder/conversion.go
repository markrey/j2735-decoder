package decoder

// #cgo CFLAGS: -I ./c/
// #cgo LDFLAGS: -L ./c/ -lasncodec
// #include <BIT_STRING.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

// octetStringToGoString takes in a ASN1 octet string and converts it to a Go string in hex
func octetStringToGoString(oString *C.OCTET_STRING_t) string {
	byteArray := *(*[]byte)(unsafe.Pointer(&oString))
	size := int(oString.size)
	str := ""
	for x := 0; x < size; x++ {
		str += fmt.Sprintf("%02X ", byteArray[x])
	}
	return strings.TrimSpace(str)
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
