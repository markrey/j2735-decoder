package decoder

// FormatType is used to identify which format to decode
type FormatType int

// Type of decoding supported by module
const (
	XML      FormatType = iota
	JSON     FormatType = iota
	SDMAPBSM FormatType = iota
	SDMAPPSM FormatType = iota
)

// ID to identify message type
const (
	BSM int64 = 20
	EVA int64 = 22
	RSA int64 = 27
	PSM int64 = 32
)

type SDMap interface {
	GetID() string
	GetHeading() int64
	SetHeading(int64)
}

// SDMapBSM contains BSM fields needed for SDMap
type SDMapBSM struct {
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

func (bsm *SDMapBSM) GetID() string {
	return bsm.ID
}

func (bsm *SDMapBSM) GetHeading() int64 {
	return bsm.Heading
}

func (bsm *SDMapBSM) SetHeading(heading int64) {
	bsm.Heading = heading
}

// SDMapPSM contains PSM fields needed for SDMap
type SDMapPSM struct {
	MsgCnt    int64
	BasicType string
	ID        string
	Lat       int64
	Long      int64
	Speed     int64
	Heading   int64
}

func (psm *SDMapPSM) GetID() string {
	return psm.ID
}

func (psm *SDMapPSM) GetHeading() int64 {
	return psm.Heading
}

func (psm *SDMapPSM) SetHeading(heading int64) {
	psm.Heading = heading
}
