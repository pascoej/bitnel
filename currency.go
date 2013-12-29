package main

type currency int

const (
	btc currency = iota
	ltc
)

func (c currency) String() string {
	switch c {
	case btc:
		return "btc"
	case ltc:
		return "ltc"
	}

	return ""
}
