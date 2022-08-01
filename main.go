package main

import (
	"fmt"

	"bookq.xyz/goods-remaining-bot/bot"
	"bookq.xyz/goods-remaining-bot/database"
	"bookq.xyz/goods-remaining-bot/oss"
)

func main() {
	database.Initialize()
	oss.Connect()
	fmt.Println("Initialize complete")
	bot.Boot()
}
