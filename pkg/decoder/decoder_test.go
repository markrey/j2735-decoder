package decoder_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

func TestBSMXML(t *testing.T) {
	file, err := os.Open("./test/bsm1.uper")
	defer file.Close()
	assert.NilError(t, err)
	xml, err := os.Open("./test/bsm1.xml")
	defer xml.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg := decoder.Decode(bytes, uint(read), decoder.XML)
	msgStr, ok := decodedMsg.(string)
	assert.Assert(t, ok)
	read, err = xml.Read(bytes)
	assert.NilError(t, err)
	xmlStr := fmt.Sprintf("%s", bytes)
	assert.Equal(t, msgStr[:802], xmlStr[:802])
}

func TestBSMJSON(t *testing.T) {
	file, err := os.Open("./test/bsm2.uper")
	defer file.Close()
	assert.NilError(t, err)
	bytes := make([]byte, 2048)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg := decoder.Decode(bytes, uint(read), decoder.JSON)
	msgStr, ok := decodedMsg.(string)
	assert.NilError(t, err)
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(msgStr), &jsonMap)
	jsonMap, ok = jsonMap["MessageFrame"].(map[string]interface{})
	assert.Assert(t, ok)
	t.Log(jsonMap)
	id, isString := jsonMap["messageId"].(string)
	assert.Assert(t, isString)
	assert.Equal(t, id, "20")
}

func TestBSMSDJSON(t *testing.T) {
	file, err := os.Open("./test/bsm1.uper")
	defer file.Close()
	assert.NilError(t, err)
	jsonF, err := os.Open("./test/bsm1sdmap.json")
	defer jsonF.Close()
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	assert.NilError(t, err)
	decodedMsg := decoder.Decode(bytes, uint(read), decoder.SDMAP)
	sdData, ok := decodedMsg.(*decoder.SDMap)
	assert.Assert(t, ok)
	read, err = jsonF.Read(bytes)
	assert.NilError(t, err)
	jsonStr := fmt.Sprintf("%s", bytes)
	t.Log(sdData)
	sdJson, mErr := json.Marshal(sdData)
	assert.NilError(t, mErr)
	assert.Equal(t, string(sdJson)[:10], jsonStr[:10])
}
