package domain

import (
	"github.com/jmoiron/sqlx"
)

type CurrencyRepositoryStub struct {
}

func (db CurrencyRepositoryStub) Get(currencyId string, from string, to string) ([]Currency, error) {

	sqlresult := []Currency{
		{
			CurrencyId: "MXN",
			UpdatedAt:  "2022-05-15 23:59:59",
			Value:      20.0123,
		},
		{
			CurrencyId: "USD",
			UpdatedAt:  "2022-05-15 23:59:59",
			Value:      1,
		},
		{
			CurrencyId: "EUR",
			UpdatedAt:  "2022-05-15 23:59:59",
			Value:      0.98,
		},
	}

	return sqlresult, nil
}

func NewCurrencyRepositoryStub(client *sqlx.DB) CurrencyRepositoryStub {
	return CurrencyRepositoryStub{}
}
