package request

import (
	"fmt"
	"reflect"
	"strconv"
)

func InterToMap(v interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if v == nil {
		return m
	}
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.String:
	case reflect.Map:
		k := value.MapRange()
		for k.Next() {
			setMap(m, k.Key().Interface(), k.Value().Interface())
		}
	case reflect.Struct:
	default:
		return m
	}
	return m
}

func setMap(m map[string]interface{}, k interface{}, v interface{}) {
	switch d := k.(type) {
	case string:
		m[d] = v
	case int, int8, int16, int32, int64:
		key := strconv.FormatInt(reflect.ValueOf(k).Int(), 10)
		m[key] = v
	case uint, uint8, uint16, uint32, uint64:
		key := strconv.FormatUint(reflect.ValueOf(k).Uint(), 10)
		m[key] = v
	case []byte:
		m[string(d)] = v
	case float32, float64:
		key := strconv.FormatFloat(reflect.ValueOf(k).Float(), 'f', -1, 64)
		m[key] = v
	case bool:
		key := strconv.FormatBool(d)
		m[key] = v
	default:
		key := fmt.Sprint(k)
		m[key] = v
	}
}

// *((*string)(unsafe.Pointer(fieldPtr))) = "admin#otokaze.cn"
