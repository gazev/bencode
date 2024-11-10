package bencode

import "io"

func Unmarshal(obj any, r io.Reader) error {
	decodedData, err := decode(r)
	if err != nil {
		return err
	}
	return populateObject(obj, decodedData)
}

func Marshal(obj any, w io.Writer) error {
	// TODO
	return nil
}
