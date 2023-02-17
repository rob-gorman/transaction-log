package utils

import (
	"encoding/json"
)

func ToJSON(payload interface{}) (js []byte, err error) {
	js, err = json.Marshal(payload)
	return js, err
}
