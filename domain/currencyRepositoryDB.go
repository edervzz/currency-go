package domain

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CurrencyRepositoryDB struct {
	client *sqlx.DB
}

func (db CurrencyRepositoryDB) Get(currencyId string, from string, to string) ([]Currency, error) {

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
	fmt.Println("query:", query)

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
	return sqlresult, nil
}

func NewCurrencyRepositoryDB(client *sqlx.DB) CurrencyRepositoryDB {
	return CurrencyRepositoryDB{
		client,
	}
}
