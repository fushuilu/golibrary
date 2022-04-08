package golibrary

import "encoding/base64"

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) (string, error) {
	if bytes, err := base64.StdEncoding.DecodeString(data); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
