package bot

import (
	"regexp"

	recordspace "bookq.xyz/goods-remaining-bot/bot/record-space"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var splitMultipleSpaces = regexp.MustCompile("\\s+")

func groupLongEvents(e Pichubot.MessageGroup) {
	for _, value := range Pichubot.LongEvents {
		if value.GroupID == e.GroupID && value.UserID == e.UserID {
			switch value.EventKey {
			case name_longEventImageInsert:
				*value.Channel <- e.Message
			}
		}
	}
}

func privateLongEvents(e Pichubot.MessagePrivate) {
	for _, value := range Pichubot.LongEvents {
		if value.UserID == e.UserID {
			switch value.EventKey {
			case name_longEventImageInsert:
				*value.Channel <- e.Message
			}
		}
	}
}

func handlerHelp(e Pichubot.MessageGroup) {
	msg := e.Message
	if msg == "/谷子bot 帮助" {
		Pichubot.SendGroupMsg(`指令太多了，写不下`, e.GroupID)
	}
}

func handlerPrivateMsgCommandParser(e Pichubot.MessagePrivate) {
	cmd_arr := splitMultipleSpaces.Split(e.Message, -1)
	if len(cmd_arr) < 2 || cmd_arr[0] != "/谷子bot" {
		return
	}
	var s string
	switch cmd_arr[1] {
	case "图库":
		s = rsCommandParser(cmd_arr[2:], e.UserID)
	case "换图":
		s = quickImageStoreChange(cmd_arr[2:], e.UserID, 0)
	default:
		s = "目前没有图库以外的操作，需要更多功能的话请向我提出请求"
	}
	if s == "" {
		return
	}
	MsgSender.Private <- QQMessage{Dst: e.UserID, S: s}
}
func handlerGroupMsgCommandParser(e Pichubot.MessageGroup) {
	var s string
	switch e.Message {
	case "/看余量":
		s = quickGetRS(recordspace.CONST_RTYPE_REMAINING, e.GroupID)
	case "/看肾表":
		s = quickGetRS(recordspace.CONST_RTYPE_BILLING, e.GroupID)
	}
	cmd_arr := splitMultipleSpaces.Split(e.Message, -1)
	if len(cmd_arr) < 2 || cmd_arr[0] != "/谷子bot" {
		return
	}
	switch cmd_arr[1] {
	case "换图":
		s = quickImageStoreChange(cmd_arr[2:], e.UserID, e.GroupID)
	}
	if s == "" {
		return
	}
	MsgSender.Group <- QQMessage{Dst: e.UserID, S: s}
}
