package cqcode

import (
	"fmt"
	"regexp"
	"strings"
)

type cqImage struct {
	All  *regexp.Regexp
	Url  *regexp.Regexp
	File *regexp.Regexp
}

var CQImage cqImage = cqImage{
	All:  regexp.MustCompile(`\[CQ:image,[0-9A-Za-z=:/?.,_-]*\]`),
	Url:  regexp.MustCompile(`url=[0-9A-Za-z=:/?._,-]+`), // conside changing back to * if + is not working
	File: regexp.MustCompile(`file=[0-9A-Za-z.]+`),
}

func (_ *cqImage) Generate(img string) string {
	return fmt.Sprintf(`[CQ:image,file=%s]`, img)
}

func (cq *cqImage) UrlTrim(img string) string {
	if len(img) < 5 {
		return ""
	}
	return strings.TrimLeft(img, "url=")
}
