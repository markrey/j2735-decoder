package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
// int xer__print2s (const void *buffer, size_t size, void *app_key)
// {
//     char *string = (char *) app_key;
//     strncat(string, buffer, size);
//     return 0;
// }
// int xer_sprint(void *string, asn_TYPE_descriptor_t *td, void *sptr)
// {
//     asn_enc_rval_t er;
//     er = xer_encode(td, sptr, XER_F_CANONICAL, xer__print2s, string);
//     if (er.encoded == -1)
//         return -1;
//     return er.encoded;
// }
import "C"
import (
	"fmt"
	"unsafe"
	"strings"
	"encoding/json"

	xj "github.com/basgys/goxml2json"
)

// FormatType is used to identify which format to decode
type FormatType int

// Type of decoding supported by module
const (
	XML  FormatType = iota
	JSON FormatType = iota
	SDBSMJSON FormatType = iota
)

// ID to identify message type
const (
	BSMID int64 = 20
)

// SDMap contains entries needed for SD Map
type SDMap struct {
	MsgCnt       int64
	ID           string
	Lat          int64
	Long         int64
	Elev         int64
	Speed        int64
	Heading      int64
	Angle        int64
}

// Decode is a public function for other packages to decode
func Decode(bytes []byte, length uint, format FormatType) string {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, uint64(length))
	if msgFrame == nil {
		Logger.Error("Cannot decode bytes to messageframe struct")
		return ""
	}
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	Logger.Infof("Decoding message type: %d", int64(msgFrame.messageId))

	size := 2048
	var buffer []byte
	for true {
		buffer = make([]byte, size)
		bufPtr := unsafe.Pointer(&buffer[0])
		rval := C.xer_sprint(bufPtr, &C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
		Logger.Infof("Bytes Encoded: %d", int(rval))
		if int(rval) == -1 {
			err := "Cannot encode message!"
			Logger.Error(err)
			return ""
		} else if int(rval) > len(buffer) {
			size = int(rval)
		}
		break
	}
	xmlStr := fmt.Sprintf("%s", buffer)
	switch format {
		case XML:
			return xmlStr
		case JSON:
			xml := strings.NewReader(xmlStr)
			json, err := xj.Convert(xml)
			if err != nil {
				Logger.Errorf("Cannot encode to JSON: %s", err)
				panic(err)
			}
			return json.String()
		case SDBSMJSON:
			if int64(msgFrame.messageId) == BSMID {
				coreData := (*C.BasicSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice)).coreData
				sdData := &SDMap{
					MsgCnt:  	int64(coreData.msgCnt),
					ID:      	octetStringToGoString(&coreData.id),
					Lat:     	int64(coreData.lat),
					Long:    	int64(coreData.Long),
					Elev:    	int64(coreData.elev),
					Speed:      int64(coreData.speed),
					Heading:    int64(coreData.heading),
					Angle:      int64(coreData.angle),
				}
				bsmJSON, _ := json.Marshal(sdData)
				Logger.Info(string(bsmJSON))
				return string(bsmJSON)
			}
			return "" 
		default:
			return xmlStr
		}	
	return ""
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
