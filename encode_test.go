package bencode

import (
	"reflect"
	"testing"
)

type File struct {
	Announce string `bencode:"announce"`
	Name     string `bencode:"name"`
	Size     int    `bencode:"size"`
}

func TestEncode1(t *testing.T) {
	f := File{
		Announce: "http://announce.com",
		Name:     "name.txt",
		Size:     10,
	}

	r, err := Marshal(f)
	if err != nil {
		t.Errorf("Marshal -> %s", err)
	}
	if !reflect.DeepEqual(r, []byte("d8:announce19:http://announce.com4:name8:name.txt4:sizei10ee")) {
		t.Errorf("invalid response")
	}
}
