package impl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// object to map[string]interface{}
func toMSI(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}
	tp := reflect.TypeOf(v)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	m := make(map[string]interface{})
	switch tp.Kind() {
	case reflect.String:
		/**
		1. json
		2. k1:v1;k2:v2
		3. k1:v1
		   k2:v2
		4. k1=v1&k2=v2
		5. k1=v1
		   k2=v2
		*/
		s := v.(string)
		if json.Valid([]byte(s)) {
			err := json.Unmarshal([]byte(s), &m)
			if err != nil {
				return nil, err
			}
			return m, nil
		}

	default:
		bts, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bts, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, nil
}

// any type to string
func anyToString(v interface{}) string {
	if v == nil {
		return ""
	}
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
