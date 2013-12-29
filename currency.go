package main

type currency struct {
	code string
	name string
}

var (
	btc = &currency{
		code: "btc",
		name: "Bitcoin",
	}

	ltc = &currency{
		code: "ltc",
		name: "Litecoin",
	}
)
