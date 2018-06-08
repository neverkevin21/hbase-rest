package rest

import (
	"encoding/base64"
)

func Base64Encode(src string) string {
	encoder := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	e := base64.NewEncoding(encoder)
	return e.EncodeToString([]byte(src))
}

func Base64Decode(src string) (string, error) {
	encoder := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	e := base64.NewEncoding(encoder)
	s, err := e.DecodeString(src)
	return string(s), err
}
