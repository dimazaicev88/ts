package tools

import (
	"encoding/json"
	"reflect"

	"github.com/guregu/null/v6"
)

func ToJson(data any) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return jsonData
}

func AsNullValue[T any](val T) null.Value[T] {
	return null.NewValue[T](val, IsNotZeroValue(val))
}

func IsNotZeroValue[T any](x T) bool {
	return !reflect.ValueOf(x).IsZero()
}
