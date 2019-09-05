package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
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
	"strings"
	"unsafe"

	xj "github.com/basgys/goxml2json"
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

// convert message frame to XML
func msgFrameToXMLString(msgFrame *C.MessageFrame_t) string {
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
			continue
		}
		break
	}
	return fmt.Sprintf("%s", buffer)
}

// xmlStringToJsonString converts xml encoded string to json
func xmlStringToJSONString(xmlStr string) string {
	xml := strings.NewReader(xmlStr)
	json, err := xj.Convert(xml)
	if err != nil {
		Logger.Errorf("Cannot encode to JSON: %s", err)
		panic(err)
	}
	return json.String()
}

// msgFrameToSDMapBSM converts message frames to format ingested by SDMAP
func msgFrameToSDMapBSM(msgFrame *C.MessageFrame_t) *SDMapBSM {
	if int64(msgFrame.messageId) != BSM {
		return nil
	}
	coreData := (*C.BasicSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice)).coreData
	sdData := &SDMapBSM {
		MsgCnt:  int64(coreData.msgCnt),
		ID:      octetStringToGoString(&coreData.id),
		Lat:     int64(coreData.lat),
		Long:    int64(coreData.Long),
		Elev:    int64(coreData.elev),
		Speed:   int64(coreData.speed),
		Heading: int64(coreData.heading),
		Angle:   int64(coreData.angle),
	}
	return sdData
}

func numToPSMType(pType int64) string {
	switch pType {
	case 0:
		return "unavailable"
	case 1:
		return "aPEDESTRIAN"
	case 2:
		return "aPEDALCYCLIST"
	case 3:
		return "aPUBLICSAFETYWORKER"
	case 4:
		return "anANIMAL"
	default:
		return "unavailable"
	}
}

// msgFrameToSDMapPSM converts message frames to format ingested by SDMAP
func msgFrametoSDMapPSM(msgFrame *C.MessageFrame_t) *SDMapPSM {
	if int64(msgFrame.messageId) != PSM {
		return nil
	}
	psmData := (*C.PersonalSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice))
	sdData := &SDMapPSM {
		MsgCnt:    int64(psmData.msgCnt),
		BasicType: numToPSMType(int64(psmData.basicType)),
		ID:        octetStringToGoString(&psmData.id),
		Lat:       int64(psmData.position.lat),
		Long:      int64(psmData.position.Long),
		Speed:     int64(psmData.speed),
		Heading:   int64(psmData.heading),
	}
	return sdData
}
