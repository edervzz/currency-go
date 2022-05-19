package domain

type Currency struct {
	CurrencyId string
	UpdatedAt  string
	Value      float32
}

type Currencies struct {
	Items []Currency
}

type CurrencyRepository interface {
	GetCurrencies(currencyId string, from string, to string) ([]Currency, error)
	InsertCurrencies(c Currencies, lastUpdatedAt string) error
	InsertLog(startTime string, timeLapse string) error
}
