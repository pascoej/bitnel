package money

import (
	"errors"
)

type Unit int64

const (
	Satoshi  Unit = 1         // .00000001
	Millibit      = 10000     // .0001
	Bitcoin       = 100000000 // 1.
)

func (u *Unit) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*u = Unit(src)
	default:
		return errors.New("money: can not find appropriate type")
	}

	return nil
}
