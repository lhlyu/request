package request

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
)

func makeSliceOfReflectValue(v reflect.Value) (slice []interface{}) {

	kind := v.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return slice
	}

	slice = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface()
	}

	return slice
}

func changeMapToURLValues(data map[string]interface{}) url.Values {
	var newUrlValues = url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			newUrlValues.Add(k, val)
		case bool:
			newUrlValues.Add(k, strconv.FormatBool(val))
		case json.Number:
			newUrlValues.Add(k, string(val))
		case int:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case int8:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case int16:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case int32:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case int64:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case uint:
			newUrlValues.Add(k, strconv.FormatUint(uint64(val), 10))
		case uint8:
			newUrlValues.Add(k, strconv.FormatUint(uint64(val), 10))
		case uint16:
			newUrlValues.Add(k, strconv.FormatUint(uint64(val), 10))
		case uint32:
			newUrlValues.Add(k, strconv.FormatUint(uint64(val), 10))
		case uint64:
			newUrlValues.Add(k, strconv.FormatUint(uint64(val), 10))
		case float64:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		case float32:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		case []byte:
			newUrlValues.Add(k, string(val))
		case []string:
			for _, element := range val {
				newUrlValues.Add(k, element)
			}
		case []int:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatInt(int64(element), 10))
			}
		case []bool:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatBool(element))
			}
		case []float64:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		case []float32:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		// these slices are used in practice like sending a struct
		case []interface{}:

			if len(val) <= 0 {
				continue
			}

			switch val[0].(type) {
			case string:
				for _, element := range val {
					newUrlValues.Add(k, element.(string))
				}
			case bool:
				for _, element := range val {
					newUrlValues.Add(k, strconv.FormatBool(element.(bool)))
				}
			case json.Number:
				for _, element := range val {
					newUrlValues.Add(k, string(element.(json.Number)))
				}
			}
		default:
			// TODO add ptr, arrays, ...
		}
	}
	return newUrlValues
}

func contains(respStatus int, statuses []int) bool {
	for _, status := range statuses {
		if status == respStatus {
			return true
		}
	}
	return false
}
