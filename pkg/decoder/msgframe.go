package decoder

// #cgo CFLAGS: -I ./c/
// #cgo LDFLAGS: -L ./c/ -lasncodec
// #include <MessageFrame.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"encoding/json"
	"fmt"
)

// MessageFrame defines the topmost structure of the different messages
type MessageFrame struct {
	MessageID int64
	Value     MessageValue
}

// MessageValue interface is used to create messageframes based on message type
type MessageValue interface {
	Create(*C.MessageFrame_t) *MessageFrame
}

// JSON converts MessageFrame into JSON format
func (msgFrame *MessageFrame) JSON() string {
	b, err := json.Marshal(msgFrame)
	if err != nil {
		panic(err)
	}
	Logger.info(fmt.Sprint(string(b)))
	return fmt.Sprint(string(b))
}
