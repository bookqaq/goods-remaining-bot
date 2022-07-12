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
	Generate func(string) string
	Trim     func(string) string
}

var CQImage cqImage = cqImage{
	All:      regexp.MustCompile(`\[CQ:image,[0-9A-Za-z=:/?.,_-]*\]`),
	Url:      regexp.MustCompile(`url=[0-9A-Za-z=:/?._,-]+`), // consider changing back to * if + is not working
	File:     regexp.MustCompile(`file=[0-9A-Za-z.]+`),
	Generate: cqImageGenerate,
	Trim:     cqImageUrlTrim,
}

func cqImageGenerate(img string) string {
	return fmt.Sprintf(`[CQ:image,file=%s]`, img)
}

func cqImageUrlTrim(img string) string {
	if len(img) < 5 {
		return ""
	}
	return strings.TrimLeft(img, "url=")
}
