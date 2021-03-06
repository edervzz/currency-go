package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CurrencyRepositoryDB struct {
	client *sqlx.DB
}

func (db CurrencyRepositoryDB) GetCurrencies(currencyId string, from string, to string) ([]Currency, error) {

	sqlresult := []Currency{}

	currencyId = strings.ToUpper(currencyId)

	query := `SELECT currency_id, updated_at, value FROM currency`

	filterCurrencyId := fmt.Sprintf("WHERE currency_id <> ''")
	if currencyId != "ALL" {
		filterCurrencyId = fmt.Sprintf("WHERE currency_id = '%s'", currencyId)
	}

	filterUpdatedAt := ""
	if from != "" && to != "" {
		filterUpdatedAt = fmt.Sprintf("AND updated_at between '%s' AND '%s'", from, to)
	} else if from != "" && to == "" {
		filterUpdatedAt = fmt.Sprintf("AND updated_at >= '%s'", from)
	} else if from == "" && to != "" {
		filterUpdatedAt = fmt.Sprintf("AND updated_at <= '%s'", to)
	}

	query = fmt.Sprintf("%s %s %s", query, filterCurrencyId, filterUpdatedAt)

	sqlRows, err := db.client.Query(query)
	if err != nil {
		return []Currency{}, err
	}

	curr := Currency{}

	for sqlRows.Next() {
		err = sqlRows.Scan(&curr.CurrencyId, &curr.UpdatedAt, &curr.Value)
		sqlresult = append(sqlresult, curr)
		if err != nil {
			return []Currency{}, err
		}
	}

	if len(sqlresult) == 0 {
		return sqlresult, errors.New("no items")
	}

	return sqlresult, nil
}

func (db CurrencyRepositoryDB) InsertCurrencies(c Currencies, lastUpdatedAt string) error {
	var atLeastOneError error
	for _, v := range c.Items {
		_, err := db.client.Exec(`INSERT INTO public.currency(currency_id, updated_at, value) VALUES ($1, $2, $3);`,
			v.CurrencyId, lastUpdatedAt, v.Value)
		if err != nil && atLeastOneError == nil {
			atLeastOneError = err
		}
	}
	return atLeastOneError
}

func (db CurrencyRepositoryDB) InsertLog(startTime string, timeLapse string) error {
	_, err := db.client.Exec(`INSERT INTO public.logger(start_datetime, time_lapse) VALUES ($1, $2);`,
		startTime, timeLapse)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func NewCurrencyRepositoryDB(client *sqlx.DB) CurrencyRepositoryDB {
	return CurrencyRepositoryDB{
		client,
	}
}
