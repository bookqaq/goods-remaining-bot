package bot

import (
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var bot = Pichubot.NewBot()

func Boot() {
	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, longEvents, handlerHelp, handlerGoodsGet, handlerGoodDelete, handlerGoodsInsert)
	bot.Config = Pichubot.Config{
		Loglvl:   Pichubot.LOGGER_LEVEL_WARNING,
		Host:     "127.0.0.1:29290",
		MasterQQ: 295589844,
		Path:     "/",
		MsgAwait: true,
	}
	bot.Run()
}
