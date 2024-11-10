package bencode

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type countWriter struct {
	*bufio.Writer
	count int
}

func (w *countWriter) WriteByte(b byte) error {
	w.count++
	return w.Writer.WriteByte(b)
}

func (w *countWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	w.count += n
	return n, err
}

func encode(obj any, w io.Writer) error {
	cw := &countWriter{
		Writer: bufio.NewWriter(w),
		count:  0,
	}
	return encodeItem(obj, cw)
}

func encodeItem(obj any, w *countWriter) error {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Int:
		return encodeInt(obj.(int), w)
	case reflect.String:
		return encodeString(obj.(string), w)
	case reflect.Slice:
		return encodeList(obj.([]any), w)
	case reflect.Map:
		// TODO
	case reflect.Struct:
		return encodeStruct(obj, w)
	default:
		return fmt.Errorf("unsupported type %s", reflect.TypeOf(obj).Kind())
	}
	return nil
}

func encodeInt(integer int, w *countWriter) error {
	if integer == 0 {
		return nil
	}
	n, err := w.Write([]byte("i"))
	if err != nil {
		return fmt.Errorf("failed writing int initial delimiter. Write('i') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing int initial delimiter -> n != 1")
	}

	str := strconv.Itoa(integer)
	n, err = w.Write([]byte(str))
	if err != nil {
		return fmt.Errorf("failed writing int value. Write(str) -> %w", err)
	}
	if n != len(str) {
		return fmt.Errorf("partial int value write")
	}

	n, err = w.Write([]byte("e"))
	if err != nil {
		return fmt.Errorf("failed writing int end delimiter. Write('e') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing int end delimiter -> n != 1")
	}
	return nil
}

func encodeString(str string, w *countWriter) error {
	if len(str) == 0 {
		return nil
	}
	lengthStr := strconv.Itoa(len(str))
	n, err := w.Write([]byte(lengthStr))
	if err != nil {
		return fmt.Errorf("failed writing string length. Write(lengthStr) -> %w", err)
	}
	if n != len(lengthStr) {
		return fmt.Errorf("partial write of string length -> n != len(lengthStr)")
	}

	n, err = w.WriteString(":" + str)
	if err != nil {
		return fmt.Errorf("failed writing string. Write(':' + str) -> %w", err)
	}
	if n != len(str)+1 {
		return fmt.Errorf("partial write of string -> n != len(str) + 1")
	}
	return nil
}

func encodeList(list []any, w *countWriter) error {
	if len(list) == 0 {
		return nil
	}
	n, err := w.Write([]byte("l"))
	if err != nil {
		return fmt.Errorf("failed writing list initial delimiter. Write('l') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing list initial delimiter -> n != 1")
	}

	for _, el := range list {
		if err := encodeItem(el, w); err != nil {
			return fmt.Errorf("failed encoding list item. encodeItem(el) -> %w", err)
		}
	}

	n, err = w.Write([]byte("e"))
	if err != nil {
		return fmt.Errorf("failed writing list initial delimiter. Write('e') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing list initial delimiter -> n != 1")
	}
	return nil
}

func encodeStruct(obj any, w *countWriter) error {
	n, err := w.Write([]byte("d"))
	if err != nil {
		return fmt.Errorf("failed writing dict initial delimiter. Write('d') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing dict start delimiter -> n != 1")
	}

	objVal := reflect.ValueOf(obj)
	for i := 0; i < objVal.NumField(); i++ {
		field := objVal.Type().Field(i)
		tag, _ := field.Tag.Lookup("bencode")
		if tag == "" {
			continue
		}
		if err := encodeString(strings.Split(tag, ",")[0], w); err != nil {
			return fmt.Errorf("failed encoding key '%s' for struct field '%s'. encodeStr() -> %w", tag, field.Name, err)
		}
		if err := encodeItem(objVal.Field(i).Interface(), w); err != nil {
			return fmt.Errorf("failed encoding key '%s' for struct field '%s'. encodeItem() -> %w", tag, field.Name, err)
		}
	}

	n, err = w.Write([]byte("e"))
	if err != nil {
		return fmt.Errorf("failed writing dict end delimiter. Write('e') -> %w", err)
	}
	if n != 1 {
		return fmt.Errorf("failed writing dict end delimiter -> n != 1")
	}
	return nil
}
