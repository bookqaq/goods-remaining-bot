package rsgroup

import "bookq.xyz/goods-remaining-bot/database"

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

func SelectRS(gp int64) ([]int32, error) {
	rows, err := database.RSGroupMapping.SelectRS.Query(gp)
	if err != nil {
		return nil, err
	}

	res := make([]int32, 0, 3)
	for rows.Next() {
		var item int32
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
