package domain

type Currency struct {
	CurrencyId string
	UpdatedAt  string
	Value      float32
}

type CurrencyRepository interface {
	Get(currencyId string, from string, to string) ([]Currency, error)
}
