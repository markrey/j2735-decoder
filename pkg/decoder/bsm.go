package decoder

// #cgo CFLAGS: -I ..
// #cgo LDFLAGS: -L .. -lasncodec
// #include <MessageFrame.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import (
	"unsafe"
)

// Accuracy defines accuracy of reading
type Accuracy struct {
	SemiMajor   int64
	SemiMinor   int64
	Orientation int64
}

// AccelSet defines the direction of acceleration
type AccelSet struct {
	Long int64
	Lat  int64
	Vert int64
	Yaw  int64
}

// Size defines the size of the vehicle
type Size struct {
	Width  int64
	Height int64
}

// Brakes define what type of brake the vehicle is using
type Brakes struct {
	WheelBrakes string
	Traction    int64
	ABS         int64
	SCS         int64
	BrakeBoost  int64
	AuxBrakes   int64
}

// BSMCoreData describes core data
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

// BasicSafetyMessage wraps the CoreData as well as extensions
type BasicSafetyMessage struct {
	CoreData *BSMCoreData
}

// Create implements the MessageValue interface to create a new BasicSafetyMessage from the C definition
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
