package main

import (
	"fmt"

	"bookq.xyz/goods-remaining-bot/bot"
	"bookq.xyz/goods-remaining-bot/database"
)

func main() {
	database.Initialize()
	fmt.Println("Initialize complete")
	bot.Boot()
}
