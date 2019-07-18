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

// Message ID values as specified in J2735 specifications
const (
	PSMID  int64 = 32
	RSAID  int64 = 27
	BSMID  int64 = 20
	SPaTID int64 = 19
	MAPID  int64 = 18
)

// Decode is a public function for other packages to decode
func Decode(bytes []byte, length uint64) *MessageFrame {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, length)
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	var message MessageValue
	Logger.Infof("Message %d decoded succesfully", int64(msgFrame.messageId))
	switch int64(msgFrame.messageId) {
	case BSMID:
		message = &BasicSafetyMessage{}
		break
	}
	return message.Create(msgFrame)
}

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
		panic(err)
	}
	return (*C.MessageFrame_t)(decoded)
}
