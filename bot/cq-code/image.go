package cqcode

import (
	"fmt"
	"regexp"
	"strings"
)

type cqImage struct {
	All,
	Url,
	File *regexp.Regexp
	Generate func(string, string, string) string
	Trim     func(string) string
}

var CQImage cqImage = cqImage{
	All:      regexp.MustCompile(`\[CQ:image,[0-9A-Za-z=:/?.,_-]*\]`),
	Url:      regexp.MustCompile(`url=[0-9A-Za-z=:/?._,-]+`), // consider changing back to * if + is not working
	File:     regexp.MustCompile(`file=[0-9A-Za-z.]+`),
	Generate: cqImageGenerate,
	Trim:     cqImageUrlTrim,
}

func cqImageGenerate(endpoint string, bucket_name string, fname string) string {
	return fmt.Sprintf(`[CQ:image,file=https://%s/%s/%s]`, endpoint, bucket_name, fname)
}

func cqImageUrlTrim(img string) string {
	if len(img) < 5 {
		return ""
	}
	return strings.TrimLeft(img, "url=")
}
