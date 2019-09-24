package decoder_test

import (
	//	"encoding/json"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

func TestBSMToXMLStringConversion(t *testing.T) {
	file, err := os.Open("./test/bsm1.uper")
	defer file.Close()
	assert.NilError(t, err)
	xml, err := os.Open("./test/bsm1.xml")
	defer xml.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg, err := decoder.DecodeString(bytes, uint(read), decoder.XML)
	assert.NilError(t, err)
	read, err = xml.Read(bytes)
	assert.NilError(t, err)
	xmlStr := fmt.Sprintf("%s", bytes)
	assert.Equal(t, decodedMsg[:802], xmlStr[:802])
}

func TestBSMToJSONStringConversion(t *testing.T) {
	file, err := os.Open("./test/bsm2.uper")
	defer file.Close()
	assert.NilError(t, err)
	bytes := make([]byte, 2048)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg, err := decoder.DecodeString(bytes, uint(read), decoder.JSON)
	assert.NilError(t, err)
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(decodedMsg), &jsonMap)
	jsonMap, ok := jsonMap["MessageFrame"].(map[string]interface{})
	assert.Assert(t, ok)
	t.Log(jsonMap)
	id, isString := jsonMap["messageId"].(string)
	assert.Assert(t, isString)
	assert.Equal(t, id, "20")
}

func TestBSMToMapAgtFormat(t *testing.T) {
	file, err := os.Open("./test/bsm1.uper")
	defer file.Close()
	assert.NilError(t, err)
	jsonF, err := os.Open("./test/bsm1map.json")
	defer jsonF.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg, err := decoder.DecodeMapAgt(bytes, uint(read), decoder.FLTBSM)
	assert.NilError(t, err)
	read, err = jsonF.Read(bytes)
	assert.NilError(t, err)
	jsonStr := fmt.Sprintf("%s", bytes)
	sdJson, mErr := json.Marshal(decodedMsg)
	assert.NilError(t, mErr)
	assert.Equal(t, string(sdJson)[:90], jsonStr[:90])
}

func TestBSMExtToMapAgtFormat(t *testing.T) {
	file, err := os.Open("./test/bsmExt2.uper")
	defer file.Close()
	assert.NilError(t, err)
	jsonF, err := os.Open("./test/bsmExt2.json")
	defer jsonF.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg, err := decoder.DecodeMapAgt(bytes, uint(read), decoder.MapAgentFormatType(0))
	assert.NilError(t, err)
	read, err = jsonF.Read(bytes)
	assert.NilError(t, err)
	jsonStr := fmt.Sprintf("%s", bytes)
	sdJson, mErr := json.Marshal(decodedMsg)
	assert.NilError(t, mErr)
	assert.Equal(t, string(sdJson)[:100], jsonStr[:100])
}

func TestPSMSDJSON(t *testing.T) {
	file, err := os.Open("./test/psm1.uper")
	defer file.Close()
	assert.NilError(t, err)
	jsonF, err := os.Open("./test/psm1map.json")
	defer jsonF.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg, err := decoder.DecodeMapAgt(bytes, uint(read), decoder.MapAgentFormatType(1))
	assert.NilError(t, err)
	read, err = jsonF.Read(bytes)
	assert.NilError(t, err)
	jsonStr := fmt.Sprintf("%s", bytes)
	sdJson, mErr := json.Marshal(decodedMsg)
	assert.NilError(t, mErr)
	assert.Equal(t, string(sdJson)[:100], jsonStr[:100])
}
