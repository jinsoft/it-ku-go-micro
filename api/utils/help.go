package utils

import "encoding/json"

func Dd(str interface{}) string {
	jsonByte, _ := json.Marshal(str)
	return string(jsonByte)
}
