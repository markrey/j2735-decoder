package main

// #cgo CFLAGS: -I ..
// #cgo LDFLAGS: -L .. -lasncodec
// #include <MessageFrame.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"flag"
	"fmt"
	"os"
	"unsafe"
)

func check(err error) {
	if err != nil {
		println(err)
		os.Exit(2)
	}
}

func main() {
	filename := flag.String("filename", "", "The file to parse")
	flag.Parse()
	if *filename == "" {
		println("Must enter a filename")
		os.Exit(2)
	}
	file, err := os.Open(*filename)
	check(err)
	bytes := make([]byte, 1024)
	read, err := file.Read(bytes)
	check(err)
	fmt.Printf("%d bytes\n", read)

	var decodedMsg unsafe.Pointer
	defer C.free(decodedMsg)
	cBytes := C.CBytes(bytes)
	defer C.free(cBytes)

	rval := C.uper_decode_complete(
		nil,
		&C.asn_DEF_MessageFrame,
		&decodedMsg,
		cBytes,
		C.ulong(read))

	msgFrame := (*C.MessageFrame_t)(decodedMsg)
	println(msgFrame.messageId)
	fmt.Printf("%d consumed\n", (uint64)(rval.consumed))
	//f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)

	// C.xer_fprint(
	// 	C.stdout,
	// 	&C.asn_DEF_MessageFrame,
	// 	decodedMsg)

	C.free_struct(
		C.asn_DEF_MessageFrame,
		decodedMsg)
	//C.uper_decode_complete(0,
}
