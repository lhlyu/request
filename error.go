package request

import (
	"errors"
	"fmt"
)

// 错误定义
const (
	// 未指定请求方法
	error_no_method_specified = "no method specified"
	// targetType 无法确定
	error_target_type_not_be_determined = "TargetType %s could not be determined"
	// 类型不支持
	error_type_not_support = "this type is not supported"
	// 当前参数仅支持字符串切片和文件流
	error_file_type_not_support = "currently only supports either a string (path/to/file), a slice of bytes (file content itself), or a os.File!"
	// 状态码值不在http包规定中
	error_status_not_exist = "StatusCode %d doesn't exist in http package"
)

func e(name, tip string, a ...interface{}) error {
	txt := name + fmt.Sprintf(tip, a...)
	return errors.New(txt)
}
