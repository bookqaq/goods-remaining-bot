package imagestore

import (
	"database/sql"
	"log"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	"bookq.xyz/goods-remaining-bot/database"
	"bookq.xyz/goods-remaining-bot/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func Exist(priv int32) (bool, error) {
	_, err := database.ImageStore.Exist.Exec(priv)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteByRS(rs int32) (int64, error) {
	rows, err := database.ImageStore.DeleteByRS.Exec(rs)
	if err != nil {
		return -1, err
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return -2, err
	}
	return affected, nil
}

func DeleteOne(priv int32) error {
	_, err := database.ImageStore.DeleteOne.Exec(priv)
	if err != nil {
		return err
	}
	return nil
}

func GetImageByRS(rs int32) ([]Image, error) {
	rows, err := database.ImageStore.SelectByRS.Query(rs)
	if err != nil {
		return nil, err
	}

	res := make([]Image, 0, 10)
	for rows.Next() {
		var item Image
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		item.RS = rs
		res = append(res, item)
	}
	return res, nil
}

func InsertImageFromMessage(msg string, rs int32) map[string]interface{} {
	var inserted []int64
	var failed []int
	data := cqcode.CQImage.All.FindAll([]byte(msg), -1)
	if data == nil {
		return nil
	}
	for i, item := range data {
		img := cqcode.CQImage.File.Find(item)
		if img == nil {
			continue
		}
		fileName := strings.TrimLeft(string(img), "file=")
		fres, err := Pichubot.GetImage(fileName)
		if err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}
		imgurl, ok := fres["data"].(map[string]interface{})["url"].(string)
		if !ok || imgurl == "" {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}

		var target Image
		target.Url, err = utils.Base64_Marshal(imgurl) // 图片有缓存，要拿真实地址
		if err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}

		target.NewPrivKey()

		res, err := database.ImageStore.InsertOne.Exec(target.Priv, target.RS, target.ImageAddBase64URL())
		if err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}
		inserted = append(inserted, id)
	}
	return map[string]interface{}{"success": inserted, "failed": failed}
}
