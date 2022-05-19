package domain

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type configuration struct {
	DBHOST     string
	DBPORT     int
	DBUSER     string
	DBPASSWORD string
	DBNAME     string
	APIKEY     string
	DEADLINE   int
	TIMEOUT    int
	ENDPOINT   string
}

func TestCurrency(t *testing.T) {
	wd, _ := os.Getwd()
	lastIdx := strings.LastIndex(wd, "/")
	wd = wd[:lastIdx]

	v := viper.New()
	v.AddConfigPath(wd + "/conf")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	c := configuration{}
	err = v.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHOST, c.DBPORT, c.DBUSER, c.DBPASSWORD, c.DBNAME)
	client, err := sqlx.Open("postgres", psqlconn)

	tn := time.Now()
	st := tn.Format("2006-01-02 15:04:05")

	db := NewCurrencyRepositoryDB(client)
	currencies := Currencies{
		Items: []Currency{
			{
				CurrencyId: "EDE",
				UpdatedAt:  st,
				Value:      99,
			},
		},
	}

	dtFrom := time.Now()

	err = db.InsertCurrencies(currencies, st)
	assert.Nil(t, err)

	dtTo := time.Now()
	timeLapse := dtTo.Sub(dtFrom)

	err = db.InsertLog(dtFrom.Format("2006-01-02 15:04:05"), timeLapse.String())
	assert.Nil(t, err)

	// err = db.InsertLog("X", timeLapse.String())
	// assert.NotNil(t, err)

	curr, err := db.GetCurrencies("ALL", "2022-05-14T23:59:59", "2022-05-16T23:59:59")
	assert.NotEmpty(t, curr)

	curr, err = db.GetCurrencies("AFN", "2022-05-14T23:59:59", "2022-05-16T23:59:59")
	assert.NotEmpty(t, curr)

	curr, err = db.GetCurrencies("AFN", "2022-05-14T23:59:59", "")
	assert.NotEmpty(t, curr)

	curr, err = db.GetCurrencies("AFN", "", "2022-05-16T23:59:59")
	assert.NotEmpty(t, curr)

	curr, err = db.GetCurrencies("AFN", "X", "2022-05-16T23:59:59")
	assert.Empty(t, curr)

	curr, err = db.GetCurrencies("AFN", "", "2021-05-16T23:59:59")
	assert.Empty(t, curr)

}
