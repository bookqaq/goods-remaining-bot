package rsgroup

import (
	"fmt"

	rsuser "bookq.xyz/goods-remaining-bot/bot/rs-user"
	"bookq.xyz/goods-remaining-bot/database"
)

type GroupMapping struct {
	ID    int32 `db:"id"`
	RS    int32 `db:"rs"`
	Group int64 `db:"gp"`
}

type ViewOpGroupGetRS struct {
	ID    int32  `db:"id"`
	Name  string `db:"name"`
	Group int64  `db:"gp"`
	RType uint8  `db:"type"`
}

func InsertOne(rs int32, gp int64) error {
	_, err := database.RSGroupMapping.InsertOne.Exec(rs, gp)
	return err
}

func SelectGP(rs int32) ([]int64, error) {
	rows, err := database.RSGroupMapping.SelectGP.Query(rs)
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0, 3)
	for rows.Next() {
		var item int64
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func SelectRS(gp int64) ([]ViewOpGroupGetRS, error) {
	rows, err := database.RSGroupMapping.SelectRS.Query(gp)
	if err != nil {
		return nil, err
	}

	res := make([]ViewOpGroupGetRS, 0, 3)
	for rows.Next() {
		var item ViewOpGroupGetRS
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func DeleteOne(rs int32, gp int64) error {
	_, err := database.RSGroupMapping.DeleteOne.Exec(rs, gp)
	return err
}

func DeleteByRS(rs int32) error {
	_, err := database.ImageStore.DeleteByRS.Exec(rs)
	return err
}

func SelectOneByRSAndGroup(rs string, group int64) (ViewOpGroupGetRS, error) {
	var res ViewOpGroupGetRS
	err := database.RSGroupMapping.SelectOneByRSAndGroup.QueryRow(group, rs).Scan(&res)
	if err != nil {
		return ViewOpGroupGetRS{}, err
	}
	return res, nil
}

func Auth(rs string, sender int64) (int32, error) {
	var rsview rsuser.ViewOpUserGetRS
	err := database.RSUserMapping.SelectOne.QueryRow(rs, sender).Scan(rsview)
	if err != nil {
		return -1, fmt.Errorf("查询图库时失败:%s", err.Error())
	}
	var rsid int32
	err = database.RecordSpace.Auth.QueryRow(rs, rsview.Owner).Scan(rsid)
	if err != nil {
		return -1, fmt.Errorf("未通过图库验证:%s", err.Error())
	}
	return rsid, nil
}
