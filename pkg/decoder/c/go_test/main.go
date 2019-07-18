package main

// #cgo CFLAGS: -I ..
// #cgo LDFLAGS: -L .. -lasncodec
// #include <MessageFrame.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func check(err error) {
	if err != nil {
		println(err)
		os.Exit(2)
	}
}

// Message ID values as specified in J2735 specifications
const (
	BSMID int64 = 20
)

// Decode is a public function for other packages to decode
func Decode(bytes []byte, length uint64) *MessageFrame {
	msgFrame := decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, length)
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	var message MessageValue
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

type MessageValue interface {
	Create(*C.MessageFrame_t) *MessageFrame
}

type Accuracy struct {
	SemiMajor   int64
	SemiMinor   int64
	Orientation int64
}

type AccelSet struct {
	Long int64
	Lat  int64
	Vert int64
	Yaw  int64
}

type Size struct {
	Width  int64
	Height int64
}

type Brakes struct {
	WheelBrakes string
	Traction    int64
	ABS         int64
	SCS         int64
	BrakeBoost  int64
	AuxBrakes   int64
}

type BSMCoreData struct {
	MsgCnt       int64
	ID           string
	SecMark      int64
	Lat          int64
	Long         int64
	Elev         int64
	Accuracy     *Accuracy
	Transmission int64
	Speed        int64
	Heading      int64
	Angle        int64
	AccelSet     *AccelSet
	Brakes       *Brakes
	Size         *Size
}

type BasicSafetyMessage struct {
	CoreData *BSMCoreData
}

func (bsm *BasicSafetyMessage) Create(msgFrame *C.MessageFrame_t) *MessageFrame {
	coreData := (*C.BasicSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice)).coreData
	return &MessageFrame{
		MessageID: BSMID,
		Value: &BasicSafetyMessage{
			CoreData: &BSMCoreData{
				MsgCnt:  int64(coreData.msgCnt),
				ID:      octetStringToGoString(&coreData.id),
				SecMark: int64(coreData.secMark),
				Lat:     int64(coreData.lat),
				Long:    int64(coreData.Long),
				Elev:    int64(coreData.elev),
				Accuracy: &Accuracy{
					SemiMajor:   int64(coreData.accuracy.semiMajor),
					SemiMinor:   int64(coreData.accuracy.semiMinor),
					Orientation: int64(coreData.accuracy.orientation),
				},
				Transmission: int64(coreData.transmission),
				Speed:        int64(coreData.speed),
				Heading:      int64(coreData.heading),
				Angle:        int64(coreData.angle),
				AccelSet: &AccelSet{
					Long: int64(coreData.accelSet.Long),
					Lat:  int64(coreData.accelSet.lat),
					Vert: int64(coreData.accelSet.vert),
					Yaw:  int64(coreData.accelSet.yaw),
				},
				Brakes: &Brakes{
					WheelBrakes: bitStringToGoString(&coreData.brakes.wheelBrakes),
					Traction:    int64(coreData.brakes.traction),
					ABS:         int64(coreData.brakes.abs),
					SCS:         int64(coreData.brakes.scs),
					BrakeBoost:  int64(coreData.brakes.brakeBoost),
					AuxBrakes:   int64(coreData.brakes.auxBrakes),
				},
				Size: &Size{
					Width:  int64(coreData.size.width),
					Height: int64(coreData.size.length),
				},
			},
		},
	}
}

type MessageFrame struct {
	MessageID int64
	Value     MessageValue
}

func (msgFrame *MessageFrame) JSON() string {
	b, err := json.Marshal(msgFrame)
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(string(b))
}

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

func main() {
	filename := flag.String("filename", "", "The file to parse")
	flag.Parse()
	if *filename == "" {
		println("Must enter a filename")
		os.Exit(2)
	}
	file, err := os.Open(*filename)
	check(err)
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	check(err)
	fmt.Printf("%d bytes\n", read)

	decodedMsg := Decode(bytes, uint64(read))
	println(decodedMsg.JSON())
	// print out BSM variables
	// fmt.Printf("-------------------BSM CORE DATA-------------------\n")
	// fmt.Printf("msgCnt, %T\n", coreData.msgCnt)
	// fmt.Printf("id, %s\n", octetStringToGoString(&coreData.id))
	// fmt.Printf("secMark, %T\n", coreData.secMark)
	// fmt.Printf("latitude, %T\n", coreData.lat)
	// fmt.Printf("longitude, %T\n", coreData.Long)
	// fmt.Printf("elevation, %T\n", coreData.elev)
	// accuracy := coreData.accuracy
	// fmt.Printf("accuracy.semiMajor, %T\n", accuracy.semiMajor)
	// fmt.Printf("accuracy.semiMinor, %T\n", accuracy.semiMinor)
	// fmt.Printf("accuracy.orientation, %T\n", accuracy.orientation)
	// fmt.Printf("transmission, %T\n", coreData.transmission)
	// fmt.Printf("speed, %T\n", coreData.speed)
	// fmt.Printf("heading, %T\n", coreData.heading)
	// fmt.Printf("angle, %T\n", coreData.angle)
	// accelSet := coreData.accelSet
	// fmt.Printf("accelSet.long, %T\n", accelSet.Long)
	// fmt.Printf("accelSet.lat, %T\n", accelSet.lat)
	// fmt.Printf("accelSet.vert, %T\n", accelSet.vert)
	// fmt.Printf("accelSet.yaw, %T\n", accelSet.yaw)
	// brakes := coreData.brakes
	// fmt.Printf("brakes.wheelBrakes, %v\n", bitStringToGoString(&brakes.wheelBrakes))
	// fmt.Printf("brakes.traction, %T\n", brakes.traction)
	// fmt.Printf("brakes.abs, %T\n", brakes.abs)
	// fmt.Printf("brakes.scs, %T\n", brakes.scs)
	// fmt.Printf("brakes.brakeBoost, %T\n", brakes.brakeBoost)
	// fmt.Printf("brakes.auxBrakes, %T\n", brakes.auxBrakes)
	// size := coreData.size
	// fmt.Printf("size.width, %T\n", size.width)
	// fmt.Printf("size.height, %T\n", size.length)
	//f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)
	// C.xer_fprint(
	// 	C.stdout,
	// 	&C.asn_DEF_MessageFrame,
	// 	decodedMsg)
	//C.uper_decode_complete(0,
}
