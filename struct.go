package bencode

import (
	"fmt"
	"reflect"
	"strings"
)

func parseTag(tag string) (string, bool) {
	if tag == "" {
		return "", true
	}
	fields := strings.Split(tag, ",")
	name := fields[0]
	omit := len(fields) > 1 && fields[1] == "omitempty"
	return name, omit
}

// I don't know how I will read this in the future, but the idea is pretty simple
func populateObject(target any, value any) error {
	ref := reflect.ValueOf(target)
	if ref.Kind() != reflect.Ptr {
		return fmt.Errorf("unexpected object type, expected ptr got %s", ref.Kind())
	}
	objVal := reflect.Indirect(ref)

	switch objVal.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			objVal.SetString(str)
			break
		}
		return fmt.Errorf("unexpected type, expected string got %s", reflect.TypeOf(value).Kind())

	case reflect.Int:
		if num, ok := value.(int); ok {
			objVal.SetInt(int64(num))
			break
		}
		return fmt.Errorf("unexpected type, expected int got %s", reflect.TypeOf(value).Kind())

	case reflect.Slice:
		if list, ok := value.([]any); ok {
			sliceVal := reflect.MakeSlice(objVal.Type(), len(list), len(list))
			for i, item := range list {
				if err := populateObject(sliceVal.Index(i).Addr().Interface(), item); err != nil {
					return fmt.Errorf("failed populating slice. populateObject() -> %w", err)
				}
			}
			objVal.Set(sliceVal)
			break
		}
		return fmt.Errorf("unexpected type, expected slice got %s", reflect.TypeOf(value).Kind())

	case reflect.Map:
		// TODO

	// We turn Bencode dictionaries into structs
	case reflect.Struct:
		dict, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("unexpected type, expected map got %s", reflect.TypeOf(value).Kind())
		}
		for i := 0; i < objVal.NumField(); i++ {
			field := objVal.Type().Field(i)
			tag, _ := field.Tag.Lookup("bencode")
			if tag == "" {
				tag = field.Name
			}

			name, omit := parseTag(tag)
			item, exists := dict[name]
			if !exists {
				if omit {
					continue
				}
				return fmt.Errorf("failed populating field '%s', key '%s' missing", field.Name, tag)
			}
			if err := populateObject(objVal.Field(i).Addr().Interface(), item); err != nil {
				return fmt.Errorf("failed populating field '%s'. populateObject() -> %w", field.Name, err)
			}
		}
	}
	return nil
}
