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

// MapAgtMsg is an interface implemented by all map agent messages
type MapAgtMsg interface {
	GetID() string
}

// MapAgtSPaT contains fields regarding multiple intersections
type MapAgtSPaT struct {
	IntersectionStateList []IntersectionState
}

// IntersectionState contains fields regarding one intersection
type IntersectionState struct {
	ID           uint64
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

// GetID gets ID of SPaT
func (spat *MapAgtSPaT) GetID() string {
	var ID = ""
	for i := 0; i < len(spat.IntersectionStateList); i++ {
		if i == len(spat.IntersectionStateList)-1 {
			ID += string(spat.IntersectionStateList[i].ID)
		} else {
			ID += string(spat.IntersectionStateList[i].ID) + ","
		}
	}
	return ID
}

// MapAgtBSM contains BSM fields needed for the map agent
type MapAgtBSM struct {
	MsgCnt  int64
	ID      string
	Lat     int64
	Long    int64
	Elev    int64
	Speed   int64
	Heading int64
	Angle   int64
	EV      int64
}

// GetID gets ID of BSM
func (bsm *MapAgtBSM) GetID() string {
	return bsm.ID
}

// MapAgtPSM contains PSM fields needed for SDMap
type MapAgtPSM struct {
	MsgCnt    int64
	BasicType string
	ID        string
	Lat       int64
	Long      int64
	Speed     int64
	Heading   int64
}

// GetID gets ID of PSM
func (psm *MapAgtPSM) GetID() string {
	return psm.ID
}
