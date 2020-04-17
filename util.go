package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

func cloneMapArray(old map[string][]string) map[string][]string {
	newMap := make(map[string][]string, len(old))
	for k, vals := range old {
		newMap[k] = make([]string, len(vals))
		for i := range vals {
			newMap[k][i] = vals[i]
		}
	}
	return newMap
}
func shallowCopyData(old map[string]interface{}) map[string]interface{} {
	if old == nil {
		return nil
	}
	newData := make(map[string]interface{})
	for k, val := range old {
		newData[k] = val
	}
	return newData
}
func shallowCopyDataSlice(old []interface{}) []interface{} {
	if old == nil {
		return nil
	}
	newData := make([]interface{}, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyFileArray(old []SuperFile) []SuperFile {
	if old == nil {
		return nil
	}
	newData := make([]SuperFile, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyCookies(old []*http.Cookie) []*http.Cookie {
	if old == nil {
		return nil
	}
	newData := make([]*http.Cookie, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}
func shallowCopyErrors(old []error) []error {
	if old == nil {
		return nil
	}
	newData := make([]error, len(old))
	for i := range old {
		newData[i] = old[i]
	}
	return newData
}

func copyRetryable(old superAgentRetryable) superAgentRetryable {
	newRetryable := old
	newRetryable.RetryableStatus = make([]int, len(old.RetryableStatus))
	for i := range old.RetryableStatus {
		newRetryable.RetryableStatus[i] = old.RetryableStatus[i]
	}
	return newRetryable
}

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
