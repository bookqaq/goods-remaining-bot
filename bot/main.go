package bot

import (
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var bot = Pichubot.NewBot()

type QQMessage struct {
	S   string
	Dst int64
}

var MsgSender struct {
	Private,
	Group chan QQMessage
}

func Boot() {
	MsgSender.Group, MsgSender.Private = make(chan QQMessage, 10), make(chan QQMessage, 10)
	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, longEvents, handlerHelp, handlerGoodsGet, handlerGoodDelete, handlerGoodsInsert)
	bot.Config = Pichubot.Config{
		Loglvl:   Pichubot.LOGGER_LEVEL_WARNING,
		Host:     "127.0.0.1:29290",
		MasterQQ: 295589844,
		Path:     "/",
		MsgAwait: true,
	}

	go func() {
		for {
			select {
			case data := <-MsgSender.Private:
				Pichubot.SendPrivateMsg(data.S, data.Dst)
			case data := <-MsgSender.Group:
				Pichubot.SendGroupMsg(data.S, data.Dst)
			}
		}
	}()

	bot.Run()
}
