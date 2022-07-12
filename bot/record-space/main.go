package recordspace

import (
	"bookq.xyz/goods-remaining-bot/database"
)

// rs.ID can be empty
func CreateOne(rs Record) error {
	_, err := database.RecordSpace.InsertOne.Exec(rs.Owner, rs.Name, rs.RType)
	return err
}

func UpdateType(targetType uint8, name string, owner int64) error {
	_, err := database.RecordSpace.UpdateType.Exec(targetType, name, owner)
	return err
}

func DeleteOne(name string, owner int64) error {
	_, err := database.RecordSpace.DeleteOne.Exec(name, owner)
	return err
}

func QueryOwnedRS(owner int64) ([]Record, error) {
	res := make([]Record, 0, 2)
	rows, err := database.RecordSpace.SelectByOwner.Query(owner)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item Record
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}
