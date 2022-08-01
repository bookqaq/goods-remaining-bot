package imagestore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	"bookq.xyz/goods-remaining-bot/database"
	"bookq.xyz/goods-remaining-bot/oss"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
	"github.com/gabriel-vasile/mimetype"
)

const go_cqhttp_relative_path = "../go-cqhttp/"

type Image struct {
	Priv  int32  `db:"priv"`
	Name  string `db:"name"`
	RS    int32  `db:"rs"`
	Fname string `db:"fname"`
}

func ImageAddBase64URL(url string) string {
	return fmt.Sprintf("%s%s", "base64://", url)
}

func (item *Image) NewPrivKey() int32 {
	for {
		key := rand.Int31()
		if err := database.ImageStore.Exist.QueryRow(key).Scan(); err == sql.ErrNoRows {
			item.Priv = key
			return key
		} else if err != nil { // if exec as expected, this case will be removed
			panic(err)
		}
	}
}

//func Exist(priv int32) (bool, error) {
//	_, err := database.ImageStore.Exist.Exec(priv)
//	if err == sql.ErrNoRows {
//		return false, nil
//	} else if err != nil {
//		return false, err
//	}
//	return true, nil
//}

func UpdateOne(id int32, fname string) error {
	_, err := database.ImageStore.UpdateOne.Exec(fname, id)
	return err
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

func SelectOne(id int32) (Image, error) {
	var res Image
	var name sql.NullString
	if err := database.ImageStore.SelectOne.QueryRow(id).Scan(&res.Priv, &res.Fname, &name, &res.RS); err != nil {
		return Image{}, err
	}
	if name.Valid {
		res.Name = name.String
	} else {
		res.Name = ""
	}
	return res, nil
}

func GetImageByRS(rs int32) ([]Image, error) {
	rows, err := database.ImageStore.SelectByRS.Query(rs)
	if err != nil {
		return nil, err
	}

	res := make([]Image, 0, 10)
	for rows.Next() {
		var item Image
		var name sql.NullString
		err = rows.Scan(&item.Priv, &item.Fname, &name)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			item.Name = name.String
		} else {
			item.Name = ""
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

		data, ok := fres["data"]
		if !ok {
			log.Printf("图片%s数据获取失败", fileName)
			failed = append(failed, i+1)
			continue
		}

		imgurl, ok := data.(map[string]interface{})["file"].(string)
		if !ok || imgurl == "" {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}

		imgurl = fmt.Sprintf("%s%s", go_cqhttp_relative_path, imgurl)
		var target Image
		target.NewPrivKey()
		mtype, err := mimetype.DetectFile(imgurl)
		if err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}
		target.Fname = fmt.Sprintf("%d%s", target.Priv, mtype.Extension())
		if err := oss.UploadFile(target.Fname, imgurl, mtype.String()); err != nil {
			log.Println(err)
			failed = append(failed, i+1)
			continue
		}

		res, err := database.ImageStore.InsertOne.Exec(target.Priv, rs, target.Fname)
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

func UpdateOneFromMessage(msg string, priv int32) error {
	img := cqcode.CQImage.File.Find([]byte(msg))
	if img == nil {
		return errors.New("未找到图片")
	}
	fileName := strings.TrimLeft(string(img), "file=")
	fres, err := Pichubot.GetImage(fileName)
	if err != nil {
		return err
	}
	data, ok := fres["data"]
	if !ok {
		return errors.New("图片拉取失败，请再试一次")
	}

	imgurl, ok := data.(map[string]interface{})["file"].(string)
	if !ok || imgurl == "" {
		return errors.New("图片路径获取，请再试一次")
	}
	imgurl = fmt.Sprintf("%s%s", go_cqhttp_relative_path, imgurl)

	var target Image
	target.NewPrivKey()
	mtype, err := mimetype.DetectFile(imgurl)
	if err != nil {
		return errors.New("图片类型判断失败")
	}
	target.Fname = fmt.Sprintf("%d%s", target.Priv, mtype.Extension())
	if err := oss.UploadFile(target.Fname, imgurl, mtype.String()); err != nil {
		return errors.New("文件上传失败")
	}
	err = UpdateOne(priv, target.Fname)
	return err
}
