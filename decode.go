package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Just to keep track of read bytes for error messages
type countReader struct {
	*bufio.Reader
	count int
}

func (r *countReader) ReadByte() (byte, error) {
	r.count++
	return r.Reader.ReadByte()
}

func decode(r io.Reader) (any, error) {
	br := &countReader{
		Reader: bufio.NewReader(r),
		count:  0,
	}

	k, err := br.ReadByte()
	if err != nil {
		return nil, err
	}
	return decodeItem(k, br)
}

func decodeItem(key byte, r *countReader) (any, error) {
	switch key {
	case 'd':
		return decodeDict(r)
	case 'l':
		return decodeList(r)
	case 'i':
		return decodeInt(r)
	default:
		if key < '0' || key > '9' {
			return nil, fmt.Errorf("invalid char %s index %d", string(key), r.count)
		}
		return decodeString(key, r)
	}
}

func decodeInt(r *countReader) (int, error) {
	integer := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("failed reading int at index %d -> %w", r.count, err)
		}
		if b == 'e' {
			break
		}
		if b < '0' || b > '9' {
			return 0, fmt.Errorf("invalid char %s reading int at index %d", string(b), r.count)
		}
		integer = integer*10 + int(b-'0')
	}
	return integer, nil
}

func decodeString(previous byte, r *countReader) (string, error) {
	length := int(previous - '0')
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", fmt.Errorf("failed decoding string length at index %d. r.ReadByte() -> %w", r.count, err)
		}
		if b == ':' {
			break
		}
		if b < '0' || b > '9' {
			return "", fmt.Errorf("invalid char %s decoding string length at index %d", string(b), r.count)
		}
		length = length*10 + int(b-'0')
	}

	read := 0
	var sb strings.Builder
	for {
		if length == read {
			break
		}
		b, err := r.ReadByte()
		if err != nil {
			return "", fmt.Errorf("failed decoding string at index %d. r.ReadByte() -> %w", r.count, err)
		}
		if err := sb.WriteByte(b); err != nil {
			return "", fmt.Errorf("failed decoding string at index %d. sb.WriteByte(%s): %w", r.count, string(b), err)
		}
		read++
	}
	return sb.String(), nil
}

func decodeList(r *countReader) ([]any, error) {
	list := make([]any, 0)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed decoding list at index %d. r.ReadByte() -> %w", r.count, err)
		}
		if b == 'e' {
			break
		}
		item, err := decodeItem(b, r)
		if err != nil {
			return nil, fmt.Errorf("failed decoding list at index %d. decodeItem() -> %w", r.count, err)
		}
		list = append(list, item)
	}
	return list, nil
}

func decodeDict(r *countReader) (map[string]any, error) {
	dict := make(map[string]any, 0)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed decoding dict at index %d. r.ReadByte() -> %w", r.count, err)
		}
		if b == 'e' {
			break
		}
		key, err := decodeString(b, r)
		if err != nil {
			return nil, fmt.Errorf("failed decoding dict at index %d. decodeString() -> %w", r.count, err)
		}
		b, err = r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed decoding dict at index %d. r.ReadByte() -> %w", r.count, err)
		}
		val, err := decodeItem(b, r)
		if err != nil {
			return nil, fmt.Errorf("failed decoding dict at index %d. decodeItem() -> %w", r.count, err)
		}
		dict[key] = val
	}
	return dict, nil
}
