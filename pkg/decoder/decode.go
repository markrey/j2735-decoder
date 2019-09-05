package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"fmt"
	"unsafe"
)

// Decode is a public function for other packages to decode
// json, xml format return full string
// sdmap format return struct
func Decode(bytes []byte, length uint, format FormatType) interface{} {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, uint64(length))
	if msgFrame == nil {
		Logger.Error("Cannot decode bytes to messageframe struct")
		return ""
	}
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	Logger.Infof("Decoding message type: %d", int64(msgFrame.messageId))

	// decode in different formats
	switch format {
		case JSON:
			return xmlStringToJSONString(msgFrameToXMLString(msgFrame))
		case SDMAPBSM:
			return msgFrameToSDMapBSM(msgFrame)
		case SDMAPPSM:
			return msgFrametoSDMapPSM(msgFrame)
		default:
			return msgFrameToXMLString(msgFrame)
	}
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
		Logger.Error(err)
		return nil
	}
	return (*C.MessageFrame_t)(decoded)
}
