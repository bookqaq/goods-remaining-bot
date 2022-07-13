package rsuser

import "bookq.xyz/goods-remaining-bot/database"

type UserMapping struct {
	ID int32 `db:"id"`
	QQ int64 `db:"dst"`
	RS int32 `db:"rs"`
}

type ViewOpUserGetRS struct {
	Owner int64  `db:"owner"`
	Name  string `db:"name"`
	RType uint8  `db:"type"`
	QQ    int64  `db:"qq"`
}

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
		if err = rows.Scan(item); err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func SelectByRS(rs int32) ([]UserMapping, error) {
	rows, err := database.RSUserMapping.SelectByRS.Query(rs)
	if err != nil {
		return nil, err
	}
	res := make([]UserMapping, 0, 10)
	for rows.Next() {
		var item UserMapping
		if err = rows.Scan(&item); err != nil {
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

// Must do auth before use
func DeleteByRS(rs int32) error {
	_, err := database.RSUserMapping.DeleteByRS.Exec(rs)
	return err
}
