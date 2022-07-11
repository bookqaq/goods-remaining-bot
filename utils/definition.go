package utils

import "fmt"

type GoodsRemainingImage struct {
	Priv int32  `db:"priv"`
	Name string `db:"name"`
	Url  string `db:"url"`
}

func (item *GoodsRemainingImage) ImageAddBase64URL() string {
	return fmt.Sprintf("%s%s", "base64://", item.Url)
}
