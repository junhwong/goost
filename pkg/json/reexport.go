package json

import jsoniter "github.com/json-iterator/go"

var API = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return API.Marshal(v)
}
func MarshalToString(v interface{}) (string, error) {
	return API.MarshalToString(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return API.Unmarshal(data, v)
}
func UnmarshalFromString(str string, v interface{}) error {
	return API.UnmarshalFromString(str, v)
}
