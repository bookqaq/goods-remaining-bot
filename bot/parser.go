package bot

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	recordspace "bookq.xyz/goods-remaining-bot/bot/record-space"
	rsuser "bookq.xyz/goods-remaining-bot/bot/rs-user"
	"bookq.xyz/goods-remaining-bot/database"
)

func rsCommandParser(commands []string, sender int64) string {
	clen := len(commands)
	if clen < 2 {
		return "未检测到有效指令"
	}

	switch {
	case commands[1] == "管理员":
		return rsUserCommandParser(commands, sender)
	case commands[1] == "群聊":
		return rsGroupCommandParser(commands, sender)
	case commands[1] == "图片":
		return rsImageStoreCommandParser(commands, sender)
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
			return fmt.Sprintf("创建链接时失败:%s", err.Error())
		}
		return fmt.Sprintf(`%s"%s"%s`, "图库", commands[1], "创建成功")
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
		err := database.RecordSpace.Auth.QueryRow(commands[1], sender).Scan(rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}

		if err = rsuser.DeleteByRS(rsid); err != nil {
			return fmt.Sprintf("删除管理员链接时失败:%s", err.Error())
		}

		if _, err = recordspace.DeleteOne(commands[1], sender); err != nil {
			return fmt.Sprintf("删除图库时失败:%s", err.Error())
		}

		return fmt.Sprintf("成功删除了图库%s", commands[1])
	case commands[0] == "查询":
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
		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}

		if clen != 2 {
			return "未进行添加操作，请检查指令格式，现阶段仅支持单次添加一个"
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
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(rsid)
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
		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}
		if clen != 2 {
			return "未进行删除操作，请检查指令格式，现阶段仅支持单次删除一个"
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
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(rsid)
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
		var rsid int32
		err := database.RecordSpace.Auth.QueryRow(rs, sender).Scan(rsid)
		if err != nil {
			return fmt.Sprintf("未通过图库验证:%s", err.Error())
		}
		if clen != 2 {
			return "未进行删除操作，请检查指令格式，现阶段仅支持单次删除一个"
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
