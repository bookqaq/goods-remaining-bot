package bot

import (
	"fmt"
	"strconv"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	rsgroup "bookq.xyz/goods-remaining-bot/bot/rs-group"
	rsuser "bookq.xyz/goods-remaining-bot/bot/rs-user"
	"bookq.xyz/goods-remaining-bot/database"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func quickImageStoreChange(commands []string, sender int64, group int64) string {
	if len(commands) < 1 {
		return "未检测到相关指令，请确认指令格式"
	}

	rss, err := rsuser.SelectRSByQQ(sender)
	if err != nil {
		return fmt.Sprintf("验证时发生错误:%s", err)
	}

	imgID, err := strconv.ParseInt(commands[0], 10, 32)
	pic_cqcode := ""
	// 图片的cq码基本都是英文，不需要考虑字符的问题
	if err != nil {
		data := cqcode.CQImage.All.FindIndex([]byte(commands[0]))
		if data == nil {
			return "指令解析失败:未找到图片"
		}
		id_tmp := commands[0][:data[0]]
		imgID, err = strconv.ParseInt(id_tmp, 10, 32)
		if err != nil {
			return fmt.Sprintf("指令解析失败:%s", err)
		}
		pic_cqcode = commands[0][data[0]:data[1]]
	}

	img, err := imagestore.SelectOne(int32(imgID))
	if err != nil || img.Priv != int32(imgID) {
		return fmt.Sprintf("验证失败:%s", err)
	}
	auth := false
	for _, rs := range rss {
		var rsid int32
		database.RecordSpace.Auth.QueryRow(rs.Name, rs.Owner).Scan(&rsid)
		if rsid == img.RS {
			auth = true
		}
	}
	if !auth {
		return "验证失败:没有修改图库的权限"
	}

	ret := "更新成功"
	if pic_cqcode == "" {
		event := Pichubot.NewEvent(sender, group, name_longEventImageUpdate)
		if group == 0 {
			MsgSender.Private <- QQMessage{Dst: sender, S: "请发送更换后的图片(发送 取消 以取消更改)"}
		} else {
			MsgSender.Group <- QQMessage{Dst: group, S: "请发送更换后的图片(发送 取消 以取消更改)"}
		}
		if err = longEventImageUpdate(event, img.Priv); err != nil {
			ret = fmt.Sprintf("更新图片失败:%s", err)
		}
	} else {
		if err = imagestore.UpdateOneFromMessage(pic_cqcode, img.Priv); err != nil {
			ret = fmt.Sprintf("更新图片失败:%s", err)
		}
	}
	return ret
}

func quickGetRS(rstype uint8, group int64) string {
	rss, err := rsgroup.SelectRS(group)
	if err != nil {
		return fmt.Sprintf("查询失败了:%s", err)
	}

	for _, rs := range rss {
		if rs.RType == rstype {
			res, err := imagestore.GetImageByRS(rs.ID)
			if err != nil {
				return fmt.Sprintf("查询图片时出错:%s", err)
			}
			if len(res) == 0 {
				return "图库中未存放图片"
			}
			var b strings.Builder
			b.WriteString("图库:")
			for _, item := range res {
				fmt.Fprintf(&b, "\n%s", cqcode.CQImage.Generate(item.Url))
			}
			return b.String()
		}
	}
	return "未找到对应的图库，请咨询管理员(或者我)"
}
