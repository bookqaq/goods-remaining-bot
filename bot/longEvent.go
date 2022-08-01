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
	for {
		msg := <-*e.Channel
		if msg == "取消" {
			return nil, true
		}
		if cqcode.CQImage.All.FindIndex([]byte(msg)) != nil {
			e.Close() // 很可能关闭channel的时候要拿锁，而这个框架没拿，导致了关了channel但没删事件的时候访问了close的channel然后panic了
			return imagestore.InsertImageFromMessage(msg, rs), false
		}
	}
}

func longEventImageUpdate(e Pichubot.LongEvent, priv int32) error {
	for {
		msg := <-*e.Channel
		if msg == "取消" {
			return errors.New("取消了添加")
		}
		if cqcode.CQImage.All.FindIndex([]byte(msg)) != nil {
			e.Close() // 同上
			return imagestore.UpdateOneFromMessage(msg, priv)
		}
	}
}
