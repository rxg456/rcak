package util

import "encoding/base64"

// base64编码
func Encode(src string) string {
	srcByte := []byte(src)
	return base64.URLEncoding.EncodeToString(srcByte)
}

// base64解码
func Decode(s string) (contents string, err error) {
	srcByte, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	contents = string(srcByte)
	return contents, err
}
