package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

func InterToHeader(v interface{}) http.Header {
	if v == nil {
		return nil
	}
	header := http.Header{}
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Map:
		switch v.(type) {
		case map[string]string:
			for key, value := range v.(map[string]string) {
				header.Set(key, value)
			}
		case map[string][]string:
			header = v.(map[string][]string)
		case http.Header:
			header = v.(http.Header)
		case map[string][]byte:
			for key, value := range v.(map[string][]byte) {
				header.Set(key, string(value))
			}
		case map[string]interface{}:
			for key, value := range v.(map[string]interface{}) {
				bts, _ := json.Marshal(value)
				header.Set(key, string(bts))
			}
		}
	case reflect.Struct:
		value := reflect.ValueOf(v)
		for i := 0; i < value.NumField(); i++ {
			name := t.Field(i).Tag.Get("head")
			if name == "" || name == "-" {
				name = underlineToConnector(t.Field(i).Name)
			}

			bts, _ := json.Marshal(value.Field(i).Interface())
			header.Set(name, string(bts))
		}
	case reflect.String:
		s := v.(string)
		sm := strToMap(s, 0)
		for key, value := range sm {
			header.Set(key, value)
		}
	}
	return header
}

// interface to map[string]string
func InterToMapStr(v interface{}) map[string]string {
	if v == nil {
		return nil
	}
	m := make(map[string]string)
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Map:
		switch v.(type) {
		case map[string]string:
			for key, value := range v.(map[string]string) {
				m[key] = value
			}
		case url.Values:
			values := v.(url.Values)
			for k, v := range values {
				m[k] = v[0]
			}
		case map[string][]byte:
			for key, value := range v.(map[string][]byte) {
				m[key] = string(value)
			}
		case map[string]interface{}:
			for key, value := range v.(map[string]interface{}) {
				m[key] = anyToString(value)
			}
		}
	case reflect.Struct:
		value := reflect.ValueOf(v)
		for i := 0; i < value.NumField(); i++ {
			name := t.Field(i).Tag.Get("json")
			if name == "" || name == "-" {
				name = toSmallTitle(t.Field(i).Name)
			}
			m[name] = anyToString(value.Field(i).Interface())
		}
	case reflect.String:
		s := v.(string)
		sm := strToMap(s, getStrType(s))
		for k, v := range sm {
			m[k] = v
		}
	}
	return m
}

// dead
func InterToMap(v interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if v == nil {
		return m
	}
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.String:
		// queryUrl

		// json

	case reflect.Map:
		k := value.MapRange()
		for k.Next() {
			setMap(m, k.Key().Interface(), k.Value().Interface())
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			name := t.Field(i).Tag.Get("json")
			if name == "" || name == "-" {
				name = toSmallTitle(t.Field(i).Name)
			}
			m[name] = value.Field(i).Interface()
		}
	default:
		return m
	}
	return m
}

// dead
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
