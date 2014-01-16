package money

// A currency defines a currency being used in the system.
type Currency int

const (
	Btc Currency = iota
	Ltc
)

func (c Currency) String() string {
	switch c {
	case Btc:
		return "btc"
	case Ltc:
		return "ltc"
	}

	return ""
}
