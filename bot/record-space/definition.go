package recordspace

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

// col in table recordSpace
type Record struct {
	ID    int32  `db:"id"`
	Owner int64  `db:"owner"`
	Name  string `db:"name"`
	RType uint8  `db:"type"` // specified value that used to simplify query command
}
