package rsuser

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
