package bot

import (
	"fmt"
	"strconv"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	recordspace "bookq.xyz/goods-remaining-bot/bot/record-space"
	rsgroup "bookq.xyz/goods-remaining-bot/bot/rs-group"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func quickImageStoreChange(commands []string, sender int64, group int64) string {
	rss, err := recordspace.QueryOwnedRS(sender)
	if err != nil {
		return fmt.Sprintf("验证失败:%s", err)
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
		if rs.ID == img.RS {
			auth = true
		}
	}
	if !auth {
		return "验证失败:没有修改图库的权限"
	}

	ret := "更新成功"
	if pic_cqcode == "" {
		event := Pichubot.NewEvent(sender, group, name_longEventImageUpdate)
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
				fmt.Fprintf(&b, "\n%s", item.Url)
			}
			return b.String()
		}
	}
	return "未找到对应的图库，请咨询管理员(或者我)"
}