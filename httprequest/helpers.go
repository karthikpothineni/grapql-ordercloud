package httprequest

import (
	"fmt"
	"reflect"
	"strings"
)

// GetValue - takes the first parameter as the actual value and second parameter as default value
// it checks if actual value is empty then assigns the default value
func GetValue(value interface{}, defaultValue interface{}) interface{} {
	if reflect.ValueOf(value).Kind() != reflect.ValueOf(defaultValue).Kind() {
		return value
	}
	switch t := value.(type) {
	case string:
		if len(strings.TrimSpace(t)) == 0 {
			return defaultValue
		}
	case uint:
		if t == 0 {
			return defaultValue
		}
	case int:
		if t == 0 {
			return defaultValue
		}
	}
	return value
}

//GetErrorString - creates a single error string
func GetErrorString(errs []error) string {
	errString := ""
	for _, err := range errs {
		errString = fmt.Sprintf("%v. %v", errString, err.Error())
	}
	return errString
}
