package rsuser

import "bookq.xyz/goods-remaining-bot/database"

func InsertOne(rs int32, dst int64) error {
	_, err := database.RSUserMapping.InsertOne.Exec(rs, dst)
	return err
}

func SelectRSByQQ(qq int64) ([]ViewOpUserGetRS, error) {
	rows, err := database.RSUserMapping.SelectRSByQQ.Query(qq)
	if err != nil {
		return nil, err
	}

	res := make([]ViewOpUserGetRS, 0, 10)
	for rows.Next() {
		var item ViewOpUserGetRS
		err = rows.Scan(item)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func DeleteOne(rs int32, dst int64) error {
	_, err := database.RSUserMapping.DeleteOne.Exec(rs, dst)
	return err
}
