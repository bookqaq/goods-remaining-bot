package rsgroup

type GroupMapping struct {
	ID    int32 `db:"id"`
	RS    int32 `db:"rs"`
	Group int64 `db:"gp"`
}

