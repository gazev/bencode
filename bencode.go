package bencode

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func Decode(obj any, r io.Reader) error {
	decodedData, err := decode(r)
	if err != nil {
		return err
	}
	return populateObject(obj, decodedData)
}

func Unmarshal(obj any, encodedStr string) error {
	decodedData, err := decode(strings.NewReader(encodedStr))
	if err != nil {
		return err
	}
	return populateObject(obj, decodedData)
}

func Encode(obj any, w io.Writer) error {
	return encode(obj, w)
}

func Marshal(obj any) ([]byte, error) {
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	if err := encode(obj, bw); err != nil {
		return nil, err
	}
	if err := bw.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
