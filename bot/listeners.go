package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	"bookq.xyz/goods-remaining-bot/database"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

const (
	name_longEventImageInsert = "goodsInsert"
)

func longEvents(e Pichubot.MessageGroup) {
	for _, value := range Pichubot.LongEvents {
		if value.GroupID == e.GroupID {
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
		Pichubot.SendGroupMsg(`可用的指令如下:
/余量帮助 显示本条指令
/添加余量 <图片> 添加余量图
/看余量 查看余量图
/删余量 <id> 删除余量图`, e.GroupID)
	}
}

func handlerGoodsGet(e Pichubot.MessageGroup) {
	if e.Message != "/看余量" {
		return
	}
	res, err := database.GoodsImages.SelectAll.Query()
	if err != nil {
		log.Println(err)
		Pichubot.SendGroupMsg(fmt.Sprintf("查询失败了(%s)", err), e.GroupID)
		return
	}
	msg := "群内余量:\n"
	for res.Next() {
		var unit imagestore.Image
		err := res.Scan(&unit.Priv, &unit.Name, &unit.Url)
		if err != nil {
			log.Println(err)
			Pichubot.SendGroupMsg(fmt.Sprintf("查询失败了(%s)", err), e.GroupID)
			return
		}
		if e.UserID == bot.Config.MasterQQ {
			msg += fmt.Sprintf("id: %d\n", unit.Priv)
		}
		msg += fmt.Sprintf("%s\n", cqcode.CQImage.Generate(unit.ImageAddBase64URL()))
	}
	msg = strings.TrimRight(msg, "\n")
	if msg == "群内余量:" {
		msg += "没有余量"
	}
	_, err = Pichubot.SendGroupMsg(msg, e.GroupID)
	if err != nil {
		log.Println(err)
	}
}

func handlerGoodDelete(e Pichubot.MessageGroup) {
	if strings.Index(e.Message, "/删余量") != 0 {
		return
	}
	splitted := strings.Split(strings.TrimRight(e.Message, "\r\n"), " ")
	if len(splitted) == 1 {
		Pichubot.SendGroupMsg("没有待删除的余量图", e.GroupID)
		return
	}
	splitted = splitted[1:]
	for _, v := range splitted {
		i, err := strconv.Atoi(v)
		if err != nil {
			Pichubot.SendGroupMsg(fmt.Sprintf("%s数字转换失败了", v), e.GroupID)
			return
		}
		_, err = database.GoodsImages.DeleteOne.Exec(i)
		if err != nil {
			log.Println(err)
			Pichubot.SendGroupMsg(fmt.Sprintf("%s删除失败了(%s)", v, err), e.GroupID)
			return
		}
	}
	Pichubot.SendGroupMsg("指定的余量图删除成功了", e.GroupID)
}

func handlerGoodsInsert(e Pichubot.MessageGroup) {
	var res map[string]interface{}

	if e.Message == "/添加余量" {
		e_long := Pichubot.NewEvent(e.Sender.UserID, e.GroupID, name_longEventImageInsert)
		Pichubot.SendGroupMsg(`请发送余量图片(输入'取消'来取消添加)`, e.GroupID)
		res_long, cancled := longEventImageInsert(e_long, rs)
		if cancled {
			Pichubot.SendGroupMsg("取消了余量添加", e.GroupID)
			return
		}
		res = res_long
	} else if strings.Index(e.Message, "/添加余量") == 0 {
		res = insertGoodsImage(e.Message)
	} else {
		return
	}

	if res == nil {
		log.Printf("消息%d-%d中未检测到图片, data=%s", e.GroupID, e.Sender.UserID, e.Message)
		return
	}

	failed, ok := res["failed"].([]int)
	if !ok {
		log.Println("失败的interface转换失败了")
		Pichubot.SendGroupMsg("失败的interface转换失败了", e.GroupID)
		return
	}
	success, ok := res["success"].([]int64)
	if !ok {
		log.Println("成功的interface转换失败了")
		Pichubot.SendGroupMsg("成功的interface转换失败了", e.GroupID)
		return
	}

	msg_ret := ""
	if lf := len(failed); lf > 0 {
		msg_ret += "第"
		for _, v := range failed {
			msg_ret += fmt.Sprintf("%d,", v)
		}
		msg_ret = strings.TrimRight(msg_ret, ",")
		msg_ret += "张图添加失败了\n"
	}
	if ls := len(success); ls > 0 {
		msg_ret += "添加成功的余量编号为:"
		for _, v := range success {
			msg_ret += fmt.Sprintf("%d,", v)
		}
		msg_ret = strings.TrimRight(msg_ret, ",")
	}

	if msg_ret == "" {
		return
	}

	_, err := Pichubot.SendGroupMsg(msg_ret, e.GroupID)
	if err != nil {
		log.Println(err)
	}
}
