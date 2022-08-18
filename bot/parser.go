package bot

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	cqcode "bookq.xyz/goods-remaining-bot/bot/cq-code"
	imagestore "bookq.xyz/goods-remaining-bot/bot/image-store"
	recordspace "bookq.xyz/goods-remaining-bot/bot/record-space"
	rsgroup "bookq.xyz/goods-remaining-bot/bot/rs-group"
	rsuser "bookq.xyz/goods-remaining-bot/bot/rs-user"
	"bookq.xyz/goods-remaining-bot/database"
	"bookq.xyz/goods-remaining-bot/oss"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func rsCommandParser(commands []string, sender int64) string {
	clen := len(commands)

	if clen == 1 && commands[0] == "查询图库" {
		res, err := rsuser.SelectRSByQQ(sender)
		if err != nil {
			if err == sql.ErrNoRows {
				return "图库:无"
			}
			return fmt.Sprintf("查询失败了:%s", err.Error())
		}

		var builder strings.Builder
		builder.Grow(2048)
		builder.WriteString("图库:")
		for _, item := range res {
			fmt.Fprintf(&builder, "\n类型:%s  名称:%s",
				recordspace.CONST_RTYPE_REVERSE_MAPPING[item.RType],
				item.Name)
		}
		return builder.String()
	}

	if clen < 2 {
		return "未检测到有效指令"
	}

	switch {
	case commands[1] == "管理员":
		return rsUserCommandParser(commands, sender)
	case commands[1] == "群聊":
		return rsGroupCommandParser(commands, sender)
	case commands[1] == "图片":
		return rsImageStoreCommandParser(commands, sender, 0)
	// record-space related
	case commands[0] == "创建图库":
		rsid, err := recordspace.CreateOne(
			recordspace.Record{
				Owner: sender,
				Name:  commands[1],
				RType: recordspace.CONST_RTYPE_DEFAULT,
			})
		if err != nil {
			return fmt.Sprintf("创建图库时失败:%s", err.Error())
		}
		if err = rsuser.InsertOne(int32(rsid), sender); err != nil {
			return fmt.Sprintf("创建链接时失败，请联系我来删除图库:%s", err.Error())
		}
		return "图库创建成功"
	case commands[1] == "修改类型":
		if clen != 3 {
			return fmt.Sprintf(`未修改图库"%s"的类型，请提供图库类型(%s\%s\%s)`,
				commands[0],
				recordspace.CONST_RTYPE_DEFAULT_STRING,
				recordspace.CONST_RTYPE_REMAINING_STRING,
				recordspace.CONST_RTYPE_BILLING_STRING,
			)
		}
		err := recordspace.TypeDetectAndUpdate(sender, commands[0], commands[2])
		if err != nil {
			return err.Error()
		}
		return "图库类型修改成功(请务必要保证每个群只有一个余量图和一个肾表)"
	case commands[0] == "删除图库":
		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(commands[1], sender).Scan(&rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}

		if err = rsuser.DeleteByRS(rsid); err != nil {
			return fmt.Sprintf("删除管理员链接时失败:%s", err.Error())
		}

		if err = rsgroup.DeleteByRS(rsid); err != nil && err != sql.ErrNoRows {
			return fmt.Sprintf("删除群聊链接时失败，请联系我进行手动删除:%s", err)
		}

		if _, err = imagestore.DeleteByRS(rsid); err != nil {
			return fmt.Sprintf("删除图库图片时失败，请联系我进行手动删除:%s", err)
		}

		if _, err = recordspace.DeleteOne(commands[1], sender); err != nil {
			return fmt.Sprintf("删除图库时失败，请联系我进行手动删除:%s", err.Error())
		}

		return fmt.Sprintf("成功删除了图库%s", commands[1])

	}
	return "未检测到相关指令，请确认指令格式"
}

func rsUserCommandParser(commands []string, sender int64) string {
	rs := commands[0]
	commands = commands[2:]
	clen := len(commands)
	if clen < 1 {
		return "未检测到相关指令，请确认指令格式"
	}

	switch commands[0] {
	case "添加":
		if clen != 2 {
			return "未进行添加操作，请检查指令格式，现阶段仅支持单次添加一个"
		}

		var rsview rsuser.ViewOpUserGetRS
		err := database.RSUserMapping.SelectOne.QueryRow(rs, sender).Scan(&rsview.Owner, &rsview.Name, &rsview.RType, &rsview.QQ)
		if err != nil {
			return fmt.Sprintf("查询图库时失败:%s", err.Error())
		}

		var rsid int32
		err = database.RecordSpace.Auth.QueryRow(rs, rsview.Owner).Scan(&rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}

		target, err := strconv.ParseInt(commands[1], 10, 64)
		if err != nil {
			return fmt.Sprintf("在转换QQ号时出现错误:%s", err.Error())
		}

		err = rsuser.InsertOne(rsid, target)
		if err != nil {
			return fmt.Sprintf("添加管理员时出现错误:%s", err.Error())
		}
		return "添加成功"
	case "查询":
		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(&rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}
		res, err := rsuser.SelectByRS(rsid)
		if err != nil {
			return fmt.Sprintf("查询管理员时出现错误:%s", err.Error())
		}
		var builder strings.Builder
		builder.Grow(2048)
		builder.WriteString("管理员(会包括图库所属者):")
		for _, item := range res {
			fmt.Fprintf(&builder, "\n%d", item.QQ)
		}
		return builder.String()
	case "删除":
		if clen != 2 {
			return "未进行删除操作，请检查指令格式，现阶段仅支持单次删除一个"
		}

		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(&rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}

		target, err := strconv.ParseInt(commands[1], 10, 64)
		if err != nil {
			return fmt.Sprintf("在转换QQ号时出现错误:%s", err.Error())
		}
		if target == sender {
			return "未进行删除操作，不能删除自身链接"
		}

		if err = rsuser.DeleteOne(rsid, target); err != nil {
			if err == sql.ErrNoRows {
				return "未找到指定管理员"
			}
			return fmt.Sprintf("删除管理员时失败:%s", err.Error())
		}
		return "删除成功"
	}
	return "未检测到相关指令，请确认指令格式"
}

func rsGroupCommandParser(commands []string, sender int64) string {
	rs := commands[0]
	commands = commands[2:]
	clen := len(commands)

	if clen < 1 {
		return "未检测到相关指令，请确认指令格式"
	}

	switch commands[0] {
	case "添加":
		if clen != 2 {
			return "未进行添加操作，请检查指令格式，现阶段仅支持单次添加一个"
		}
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}
		target, err := strconv.ParseInt(commands[1], 10, 64)
		if err != nil {
			return fmt.Sprintf("在转换群号时出现错误:%s", err.Error())
		}

		err = rsgroup.InsertOne(rsid, target)
		if err != nil {
			return fmt.Sprintf("添加群聊时出现错误:%s", err.Error())
		}
		return "添加成功"
	case "查询":
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}
		res, err := rsgroup.SelectGP(rsid)
		if err != nil {
			return fmt.Sprintf("查询群聊时出现错误:%s", err.Error())
		}
		var builder strings.Builder
		builder.Grow(2048)
		builder.WriteString("绑定的群聊:")
		for _, item := range res {
			fmt.Fprintf(&builder, "\n%d", item)
		}
		return builder.String()
	case "删除":
		if clen != 2 {
			return "未进行删除操作，请检查指令格式，现阶段仅支持单次删除一个"
		}
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}

		target, err := strconv.ParseInt(commands[1], 10, 64)
		if err != nil {
			return fmt.Sprintf("在转换群号时出现错误:%s", err.Error())
		}

		if err = rsgroup.DeleteOne(rsid, target); err != nil {
			if err == sql.ErrNoRows {
				return "未找到指定群聊"
			}
			return fmt.Sprintf("删除群聊链接时失败:%s", err.Error())
		}
		return "删除成功"
	}
	return "未检测到相关指令，请确认指令格式"
}

func rsImageStoreCommandParser(commands []string, sender int64, group int64) string {
	rs := commands[0]
	commands = commands[2:]
	clen := len(commands)

	if clen < 1 {
		return "未检测到相关指令，请确认指令格式"
	}

	switch {
	case commands[0] == "添加":
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}

		event := Pichubot.NewEvent(sender, group, name_longEventImageInsert)
		if group == 0 {
			MsgSender.Private <- QQMessage{Dst: sender, S: `请发送需要添加的图片(发送 取消 以取消添加)`}
		}
		res, cancled := longEventImageInsert(event, rsid)
		if cancled {
			return "取消了图片添加"
		}
		return imagestore.ParseImageInsertResult(res)

	case strings.Index(commands[0], "添加") == 0:
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}

		res := imagestore.InsertImageFromMessage(commands[0], rsid)
		return imagestore.ParseImageInsertResult(res)

	case commands[0] == "查询" && group == 0:
		rsid, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return fmt.Sprintf("查询图库时出错:%s", err)
		}
		res, err := imagestore.GetImageByRS(rsid)
		if err != nil {
			return fmt.Sprintf("查询图片时出错:%s", err)
		}
		if len(res) == 0 {
			return "图库中未存放图片"
		}
		var b strings.Builder
		b.WriteString("图库:")
		for _, item := range res {
			fmt.Fprintf(&b, "\nid:%d\n%s", item.Priv, cqcode.CQImage.Generate(oss.Endpoint, oss.Bucket_name, item.Fname))
		}
		return b.String()
	case commands[0] == "查询":
		rs, err := rsgroup.SelectOneByRSAndGroup(rs, group)
		if err != nil {
			return ""
		}
		res, err := imagestore.GetImageByRS(rs.ID)
		if err != nil {
			return ""
		}
		if len(res) == 0 {
			return ""
		}
		var b strings.Builder
		b.WriteString("图库:")
		for _, item := range res {
			fmt.Fprintf(&b, "\n%s", cqcode.CQImage.Generate(oss.Endpoint, oss.Bucket_name, item.Fname))
		}
		return b.String()
	case commands[0] == "删除":
		if clen != 2 {
			return "未进行删除操作，请检查指令格式，现阶段仅支持单次删除一个"
		}
		_, err := rsgroup.Auth(rs, sender)
		if err != nil {
			return err.Error()
		}

		imgid, err := strconv.ParseInt(commands[1], 10, 32)
		if err != nil {
			return err.Error()
		}

		if err := imagestore.DeleteOne(int32(imgid)); err != nil {
			return fmt.Sprintf("删除图片时失败:%s", err)
		}
		return "删除成功"
	}
	return "未检测到相关指令，请确认指令格式"
}
