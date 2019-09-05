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
}

// SDMapPSM contains PSM fields needed for SDMap
type SDMapPSM struct {
	MsgCnt    int64
	BasicType string
	ID        string
	Lat       int64
	Long      int64
//	Elev      int64
	Speed     int64
	Heading   int64
}
