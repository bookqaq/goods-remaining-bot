package bot

import (
	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
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
