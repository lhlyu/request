package request

import (
	"fmt"
	"reflect"
	"strconv"
)

// any type to string
func anyToString(v interface{}) string {
	switch d := v.(type) {
	case string:
		return d
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case []byte:
		return string(d)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(d)
	default:
		return fmt.Sprint(v)
	}
}

func toMSS(v interface{}) MSS{
	tp := reflect.TypeOf(v)
	switch tp.Kind() {
	case reflect.Map:
		switch data := v.(type) {
		case MSS:
			return data
		case MSI:
			m := make(MSS)
			for k, v := range data {
				m[k] = anyToString(v)
			}
			return m
		}
	case reflect.String:
		s := v.(string)
		data := strToMSI(s,getStrType(s))
		m := make(MSS)
		for k, v := range data {
			m[k] = anyToString(v)
		}
		return m
	}
	return nil
}

func toMSI(v interface{}) MSI{
	tp := reflect.TypeOf(v)
	switch tp.Kind() {
	case reflect.Map:
		switch data := v.(type) {
		case MSI:
			return data
		}
	case reflect.String:
		s := v.(string)
		return  strToMSI(s,getStrType(s))
	}
	return nil
}
