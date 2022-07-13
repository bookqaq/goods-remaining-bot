package recordspace

import "bookq.xyz/goods-remaining-bot/database"

const (
	// default RType, no special meaning
	CONST_RTYPE_DEFAULT = iota
	// for image set of remaining goods
	CONST_RTYPE_REMAINING
	// for image set of billings
	CONST_RTYPE_BILLING
	// reserved, to be expanded
	CONST_RTYPE_OTHERS
)

const (
	CONST_RTYPE_DEFAULT_STRING   = "无"
	CONST_RTYPE_REMAINING_STRING = "余量图"
	CONST_RTYPE_BILLING_STRING   = "肾表"
	CONST_RTYPE_OTHERS_STRING    = "其他"
)

var CONST_RTYPE_REVERSE_MAPPING = [4]string{"无", "余量图", "肾表", "其他"}

// col in table recordSpace
type Record struct {
	ID    int32  `db:"id"`
	Owner int64  `db:"owner"`
	Name  string `db:"name"`
	RType uint8  `db:"type"` // specified value that used to simplify query command
}

// rs.ID can be empty
func CreateOne(rs Record) (int64, error) {
	rows, err := database.RecordSpace.InsertOne.Exec(rs.Owner, rs.Name, rs.RType)
	if err != nil {
		return -1, err
	}
	rsid, err := rows.LastInsertId()
	return rsid, err
}

func UpdateType(targetType uint8, name string, owner int64) (int64, error) {
	rows, err := database.RecordSpace.UpdateType.Exec(targetType, name, owner)
	if err != nil {
		return -1, err
	}
	affected, err := rows.RowsAffected()
	return affected, err
}

func DeleteOne(name string, owner int64) (int64, error) {
	rows, err := database.RecordSpace.DeleteOne.Exec(name, owner)
	if err != nil {
		return -1, err
	}
	affected, err := rows.RowsAffected()
	return affected, err
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
