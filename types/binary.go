package types

import (
	"bytes"
	"encoding/binary"
)

func EncodeBytesToBinary(byteFile []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, byteFile)
	//log.Println(buf.Bytes())
	return buf.Bytes(), err
}
