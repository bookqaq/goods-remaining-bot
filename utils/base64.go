package utils

import (
	"encoding/base64"
	"io"
	"net/http"
)

// Fetch image from url and encode it to base64
func Base64_Marshal(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	img_b64 := base64.StdEncoding.EncodeToString(img)
	return img_b64, nil
}

// !!! 传统写法cap和len对不上 !!!
// Mrs4s/go-cqhttp internal/param/param.go#L90
// Base64DecodeString decode base64 with avx2
// see https://github.com/segmentio/asm/issues/50
// avoid incorrect unsafe usage in origin library
//func Base64DecodeString(s string) ([]byte, error) {
//	e := base64.StdEncoding
//	dst := make([]byte, e.DecodedLen(len(s)))
//	n, err := e.Decode(dst, utils.S2B(s))
//	return dst[:n], err
//}
