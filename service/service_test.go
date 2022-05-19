package service

import (
	"currency-go/domain"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name    string
	request CurrencyRequest
}
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

func TestCurrencyGet(t *testing.T) {
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

	s := NewDefaultCurrencyService(domain.NewCurrencyRepositoryDB(client))

	tc := []testCase{
		{
			name: "OK: All + From + To",
			request: CurrencyRequest{
				CurrencyId: "ALL",
				From:       "2022-05-14T23:59:59",
				To:         "2022-05-16T23:59:59",
			},
		},
		{
			name: "OK: All + From",
			request: CurrencyRequest{
				CurrencyId: "ALL",
				From:       "2022-05-14T23:59:59",
				To:         "",
			},
		},
		{
			name: "OK: All + To",
			request: CurrencyRequest{
				CurrencyId: "ALL",
				From:       "",
				To:         "2022-05-16T23:59:59",
			},
		},
		{
			name: "ERR: Not found",
			request: CurrencyRequest{
				CurrencyId: "ZZZ",
				From:       "2022-05-14T23:59:59",
				To:         "2022-05-15T23:59:59",
			},
		},
		{
			name: "ERR: FROM < TO",
			request: CurrencyRequest{
				CurrencyId: "ZZZ",
				From:       "2022-05-15T23:59:59",
				To:         "2022-05-14T23:59:59",
			},
		},
		{
			name: "ERR: bad config",
			request: CurrencyRequest{
				CurrencyId: "ZZZ",
				From:       "2022-05-15T23:59:59",
				To:         "2022-05-14T23:59:59",
			},
		},
	}

	for _, v := range tc {
		res, appErr := s.Get(v.request)
		if strings.Contains(v.name, "OK:") {
			assert.NotEmpty(t, res)
			assert.Empty(t, appErr)
		} else if strings.Contains(v.name, "ERR:") {
			assert.Equal(t, 0, len(res.Items))
			assert.NotEqual(t, 200, appErr.Code)
		}
	}
	client.Close()
}
