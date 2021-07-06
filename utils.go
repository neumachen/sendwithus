package sendwithus

import (
	"reflect"
	"strings"
)

// IsNil checks if the given interface is nil.
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	valueOf := reflect.ValueOf(v)

	switch valueOf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return valueOf.IsNil()
	}

	return false
}

// OneIsNil ...
func OneIsNil(v ...interface{}) bool {
	for i := range v {
		if IsNil(v[i]) {
			return true
		}
	}
	return false
}

// StringIsEmpty checks if the givne str is empty. This does not consider
// white spaces/tab as non empty values
func StringIsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
