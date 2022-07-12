package imagestore

import (
	"database/sql"
	"fmt"
	"math/rand"

	"bookq.xyz/goods-remaining-bot/database"
)

type Image struct {
	Priv int32  `db:"priv"`
	Name string `db:"name"`
	Url  string `db:"url"`
	RS   int32  `db:"rs"`
}

func (item *Image) ImageAddBase64URL() string {
	return fmt.Sprintf("%s%s", "base64://", item.Url)
}

func (item *Image) NewPrivKey() int32 {
	for {
		key := rand.Int31()
		if _, err := database.ImageStore.Exist.Query(key); err == sql.ErrNoRows {
			item.Priv = key
			return key
		} else if err != nil { // if exec as expected, this case will be removed
			panic(err)
		}
	}
}
