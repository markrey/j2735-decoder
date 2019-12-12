package decoder

// StringFormatType is used to identify which format to decode
type StringFormatType int

// MapAgentFormatType is used to identify
type MapAgentFormatType int

// String formats that the module supports
const (
	XML  StringFormatType = iota
	JSON StringFormatType = iota
)

// Map Agent format type that the module supports
// These are flat json maps
const (
	FLTBSM  MapAgentFormatType = iota
	FLTPSM  MapAgentFormatType = iota
	MAPSPaT MapAgentFormatType = iota
)

// ID to identify message type
const (
	SPaT int64 = 19
	BSM  int64 = 20
	EVA  int64 = 22
	RSA  int64 = 27
	PSM  int64 = 32
)

// MapAgtMsg is interface for all map agt messages
type MapAgtMsg interface {
	GetID() string
	SetTopic(topic string)
}

// MapAgtPlugin is an mixin for all map agt messages
type MapAgtPlugin struct {
	ID    string
	Topic string
}

// GetID gets ID of map agent entity
func (agt *MapAgtPlugin) GetID() string {
	return agt.ID
}

// SetTopic sets topic of map agent entity
func (agt *MapAgtPlugin) SetTopic(topic string) {
	agt.Topic = topic
}

// SPaTList contains fields regarding multiple intersections
type SPaTList struct {
	IntersectionStateList []IntersectionState
}

// IntersectionState contains fields regarding one intersection
type IntersectionState struct {
	MapAgtPlugin
	MinuteOfYear uint64
	TimeStamp    uint64
	SignalPhases []SignalPhaseGroup
}

// SignalPhaseGroup containers fields for one signal group
type SignalPhaseGroup struct {
	GroupID    uint64
	Status     uint64
	MaxEndTime uint64
	MinEndTime uint64
}

// MapAgtBSM contains BSM fields needed for the map agent
type MapAgtBSM struct {
	MapAgtPlugin
	MsgCnt  int64
	Lat     int64
	Long    int64
	Elev    int64
	Speed   int64
	Heading int64
	Angle   int64
	EV      int64
}

// MapAgtPSM contains PSM fields needed for SDMap
type MapAgtPSM struct {
	MapAgtPlugin
	MsgCnt    int64
	BasicType string
	Lat       int64
	Long      int64
	Speed     int64
	Heading   int64
}
