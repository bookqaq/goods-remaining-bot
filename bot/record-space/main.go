package recordspace

import (
	"errors"
	"fmt"
)

func TypeDetectAndUpdate(owner int64, rsname, rstype string) error {
	var affected int64
	var err error
	switch rstype {
	case CONST_RTYPE_DEFAULT_STRING:
		affected, err = UpdateType(CONST_RTYPE_DEFAULT, rsname, owner)
	case CONST_RTYPE_REMAINING_STRING:
		affected, err = UpdateType(CONST_RTYPE_REMAINING, rsname, owner)
	case CONST_RTYPE_BILLING_STRING:
		affected, err = UpdateType(CONST_RTYPE_BILLING, rsname, owner)
	default:
		affected, err = UpdateType(CONST_RTYPE_OTHERS, rsname, owner)
	}
	if err != nil {
		return err
	}
	if affected < 1 {
		return errors.New(fmt.Sprintf("%s:%s", "未找到图库", rsname))
	}
	return nil
}
