package main

import (
	"bookq.xyz/goods-remaining-bot/bot"
	"bookq.xyz/goods-remaining-bot/database"
)

func main() {
	database.Initialize()
	bot.Boot()
}
