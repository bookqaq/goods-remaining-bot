package imagestore

import (
	"fmt"
	"strings"
)

func ParseImageInsertResult(res map[string]interface{}) string {
	if res == nil {
		return "消息中未检测到图片"
	}

	failed, ok := res["failed"].([]int)
	if !ok {
		return `convert res["failed"] fails`
	}
	success, ok := res["success"].([]int64)
	if !ok {
		return `convert res["success"] fails`
	}

	var b strings.Builder
	b.Grow(1024)
	if lf := len(failed); lf > 0 {
		b.WriteString("第")
		for _, v := range failed {
			fmt.Fprintf(&b, "%d ", v)
		}
		b.WriteString("张图片添加失败了\n")
	}
	if ls := len(success); ls > 0 {
		b.WriteString("添加成功的编号为:")
		for _, v := range success {
			fmt.Fprintf(&b, "%d ", v)
		}
	}
	return b.String()
}
