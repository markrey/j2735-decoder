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
	FLTBSM MapAgentFormatType = iota
	FLTPSM MapAgentFormatType = iota
)

// ID to identify message type
const (
	BSM int64 = 20
	EVA int64 = 22
	RSA int64 = 27
	PSM int64 = 32
)

// MapAgtMsg is an interface implemented by all map agent messages
type MapAgtMsg interface {
	GetID() string
	GetHeading() int64
	SetHeading(int64)
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

// GetHeading gets heading of BSM
func (bsm *MapAgtBSM) GetHeading() int64 {
	return bsm.Heading
}

// SetHeading sets heading of BSM
func (bsm *MapAgtBSM) SetHeading(heading int64) {
	bsm.Heading = heading
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

// GetHeading gets heading of PSM
func (psm *MapAgtPSM) GetHeading() int64 {
	return psm.Heading
}

// SetHeading sets heading of PSM
func (psm *MapAgtPSM) SetHeading(heading int64) {
	psm.Heading = heading
}
