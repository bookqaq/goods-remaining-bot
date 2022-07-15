package bot

import (
	"errors"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

const (
	name_longEventImageInsert = "IMG_INSERT"
	name_longEventImageUpdate = "IMG_UPDATE"
)

func longEventImageInsert(e Pichubot.LongEvent, rs int32) (map[string]interface{}, bool) {
	defer e.Close()
	for {
		msg := <-*e.Channel
		if msg == "取消" {
			return nil, true
		}
		if cqcode.CQImage.All.FindIndex([]byte(msg)) != nil {
			return imagestore.InsertImageFromMessage(msg, rs), false
		}
	}
}

func longEventImageUpdate(e Pichubot.LongEvent, priv int32) error {
	defer e.Close()
	for {
		msg := <-*e.Channel
		if msg == "取消" {
			return errors.New("取消了添加")
		}
		if cqcode.CQImage.All.FindIndex([]byte(msg)) != nil {
			return imagestore.UpdateOneFromMessage(msg, priv)
		}
	}
}
