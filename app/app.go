package app

import (
	"context"
	"currency-go/domain"
	"currency-go/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var client *sqlx.DB

func Run() {
	wd, _ := os.Getwd()
	config := PrepareConfiguration(wd + "/conf")

	client := NewClientDB(config)
	defer client.Close()

	go do(config, client)

	ch := NewCurrencyHandler(service.NewDefaultCurrencyService(domain.NewCurrencyRepositoryDB(client)))

	router := mux.NewRouter()
	router.HandleFunc("/currencies/{id:[a-z]{3}}", ch.Get).Methods(http.MethodGet)
	// router.HandleFunc("/currencyApiTest", ch.GetCurrencyAPISimulate).Methods(http.MethodGet)
	fmt.Println("listening on : 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		fmt.Println(err.Error())
	}

}

func NewClientDB(c configuration) *sqlx.DB {
	var err error
	if client == nil {
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHOST, c.DBPORT, c.DBUSER, c.DBPASSWORD, c.DBNAME)
		client, err = sqlx.Open("postgres", psqlconn)
		if err != nil {
			panic(err)
		}
		err = client.Ping()
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected!")
	}

	return client
}

func do(c configuration, client *sqlx.DB) {
	parentCtx := context.Background()
	lauchWorker := make(chan bool, 1)
	workerId := make(chan int, 1)
	childCtx, cancelDeadline := context.WithDeadline(parentCtx, getDeadline(c.DEADLINE))
	defer cancelDeadline()

	lauchWorker <- true

	for newWorker := 1; ; newWorker++ {
		select {
		case <-childCtx.Done():
			fmt.Println("RESET WORKER ================================================================")
			fmt.Println("DEADLINE:", c.DEADLINE, "| TIMEOUT:", c.TIMEOUT, "| ENDPOINT:", c.ENDPOINT)
			childCtx, cancelDeadline = context.WithDeadline(parentCtx, getDeadline(c.DEADLINE))
			lauchWorker <- true
			break
		case <-lauchWorker:
			fmt.Println("START WORKER :", time.Now())
			workerCtx, cancelWorkerCtx := context.WithTimeout(childCtx, getTimeout(c.TIMEOUT))
			go worker(workerCtx, c, client, workerId)
			result := <-workerId
			fmt.Println("END WORKER   :", time.Now())
			cancelWorkerCtx()
			if result == 1 {
				break
			}
		}
	}
}

func worker(ctx context.Context, config configuration, client *sqlx.DB, workerId chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("2.1 cancelled: WID:", workerId)
			workerId <- 1
			return
		default:
			id := 0
			dtFrom := time.Now()
			response, err := getCurrenciesFromAPI(ctx, config.ENDPOINT)
			if err != nil {
				fmt.Println("error:", err.Error())
				id = 1
				workerId <- id
				return
			}
			dtTo := time.Now()
			db := domain.NewCurrencyRepositoryDB(client)

			currencies := domain.Currencies{}
			currItem := domain.Currency{}

			for _, v := range response.Data {
				currItem.CurrencyId = v.Code
				currItem.Value = v.Value
				currencies.Items = append(currencies.Items, currItem)
			}

			err = db.InsertCurrencies(currencies, response.Meta.LastUpdatedAt)
			if err == nil {
				timeLapse := dtTo.Sub(dtFrom)
				db.InsertLog(dtFrom.String(), timeLapse.String())
			}

			workerId <- id
			return
		}
	}
}

func getCurrenciesFromAPI(ctx context.Context, endpoint string) (*service.CurrencyPayloadAPI, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		return &service.CurrencyPayloadAPI{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if e, ok := err.(net.Error); ok && e.Timeout() {
		log.Printf("Do request timeout: %s\n", err)
		return &service.CurrencyPayloadAPI{}, err
	} else if err != nil {
		log.Printf("Cannot do request: %s\n", err)
		return &service.CurrencyPayloadAPI{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return &service.CurrencyPayloadAPI{}, err
	}

	var responseObject service.CurrencyPayloadAPI
	json.Unmarshal(bodyBytes, &responseObject)
	return &responseObject, nil
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

func PrepareConfiguration(path string) configuration {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	config := configuration{}
	err := v.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return config
}

func getDeadline(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}

func getTimeout(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
