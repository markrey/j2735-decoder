package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yh742/j2735-decoder/pkg/decoder"
)

func main() {
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		os.Exit(5)
	}
	reader := bufio.NewReader(file)
	lineCnt := 0
	for true {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			break
		}
		splits := strings.Split(line, ":")
		hexString := strings.TrimSpace(splits[len(splits)-1])
		data, err := hex.DecodeString(hexString)
		if err != nil {
			continue
		}
		decodedMsg, err := decoder.DecodeMapAgt(data,
			uint(len(data)),
			decoder.MapAgentFormatType(0))
		if err != nil {
			println(err)
		}
		b, err := json.Marshal(decodedMsg)
		fmt.Printf("%s\n", string(b))
		lineCnt++
	}
}
