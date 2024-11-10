package bencode

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecode1(t *testing.T) {
	fmt.Println("Runnig Test1")
	r := strings.NewReader("d8:announcei2ee")
	data, err := decode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	mapData, ok := data.(map[string]any)
	if !ok {
		t.Errorf("decode returned an invalid type")
	}

	announce, exists := mapData["announce"]
	if !exists {
		t.Errorf("missing 'announce' key")
	}

	val, ok := announce.(int)
	if !ok {
		t.Errorf("invalid type, expected int")
	}

	if val != 2 {
		t.Errorf("wrong announce value")
	}
}

func TestDecode2(t *testing.T) {
	fmt.Println("Runnig Test2")
	r := strings.NewReader("i2e")
	data, err := decode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	val, ok := data.(int)
	if !ok {
		t.Errorf("decode returned an invalid type")
	}

	if val != 2 {
		t.Errorf("wrong announce value")
	}
}

func TestDecode3(t *testing.T) {
	fmt.Println("Runnig Test3")
	r := strings.NewReader("l3:egg4:foose")
	data, err := decode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	val, ok := data.([]any)
	if !ok {
		t.Errorf("decode returned an invalid type")
	}

	str, ok := val[0].(string)
	if !ok {
		t.Errorf("decode returned an invalid type")
	}

	if str != "egg" {
		t.Errorf("wrong key at index 0")
	}

	str, ok = val[1].(string)
	if !ok {
		t.Errorf("decode returned an invalid type")
	}

	if str != "foos" {
		t.Errorf("wrong key at index 0")
	}
}
