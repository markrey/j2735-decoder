package decoder

// #cgo CFLAGS: -I ./c/
// #cgo LDFLAGS: -L ./c/ -lasncodec
// #include <MessageFrame.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	PSM_ID  int = 32
	RSA_ID  int = 27
	BSM_ID  int = 20
	SPaT_ID int = 19
	MAP_ID  int = 18
)

// decodeMessageFrame requires caller to free the MessageFrame returned
func decodeMessageFrame(descriptor *C.asn_TYPE_descriptor_t, bytes []byte, length uint64) *C.MessageFrame_t {
	var decoded unsafe.Pointer
	cBytes := C.CBytes(bytes)
	defer C.free(cBytes)
	rval := C.uper_decode_complete(
		nil,
		descriptor,
		&decoded,
		cBytes,
		C.ulong(length))
	if rval.code != C.RC_OK {
		err := fmt.Sprintf("Broken Rectangle encoding at byte %d", (uint64)(rval.consumed))
		Logger.Error(err)
		panic(err)
	}
	return (*C.MessageFrame_t)(decoded)
}
