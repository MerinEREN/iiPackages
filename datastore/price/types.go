package price

type Price struct {
	Amount   float64 `json:amount"`
	Currency string  `json:currency"`
}
