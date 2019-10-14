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
	"errors"
	"fmt"
	"unsafe"
)

// DecodeString is a public function for other packages to decode string
// return: string in either json, xml format
func DecodeString(bytes []byte, length uint, format StringFormatType) (string, error) {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, uint64(length))
	if msgFrame == nil {
		Logger.Error("Cannot decode bytes to messageframe struct")
		return "", errors.New("Cannot decode bytes to messageframe struct")
	}
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	Logger.Infof("Decoding message type: %d", int64(msgFrame.messageId))
	// decode in different formats
	switch format {
	case JSON:
		xml, err := msgFrameToXMLString(msgFrame)
		if err != nil {
			return "", errors.New("decoding xml error")
		}
		return xmlStringToJSONString(xml)
	case XML:
		return msgFrameToXMLString(msgFrame)
	default:
		return "", errors.New("format type not supported")
	}
}

// DecodeMapAgt is a public function for struct types used for SD Maps
// return: SDMapAgt interface
func DecodeMapAgt(bytes []byte, length uint, format MapAgentFormatType) (MapAgtMsg, error) {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, uint64(length))
	if msgFrame == nil {
		Logger.Error("Cannot decode bytes to messageframe struct")
		return nil, errors.New("Cannot decode bytes to messageframe struct")
	}
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	Logger.Infof("Decoding message type: %+v", format)
	switch format {
	case FLTBSM:
		return msgFrameToSDMapBSM(msgFrame)
	case FLTPSM:
		return msgFrametoSDMapPSM(msgFrame)
	case MAPSPaT:
		return msgFrametoMapSPaT(msgFrame)
	default:
		return nil, errors.New("format type not supported")
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
