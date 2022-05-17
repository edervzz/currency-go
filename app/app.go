package app

import (
	"context"
	"currency-go/domain"
	"currency-go/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "currency"
)

var client *sqlx.DB

func Run() {
	go do()

	client := NewClientDB()
	defer client.Close()

	ch := NewCurrencyHandler(service.NewDefaultCurrencyService(domain.NewCurrencyRepositoryDB(client)))

	router := mux.NewRouter()
	router.HandleFunc("/currencies/{id}", ch.Get).Methods(http.MethodGet)
	router.HandleFunc("/currencyApiTest", ch.GetCurrencyAPISimulate).Methods(http.MethodGet)
	fmt.Println("listening on : 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		fmt.Println(err.Error())
	}

}

func NewClientDB() *sqlx.DB {
	var err error
	if client == nil {
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

func do() {
	parentCtx := context.Background()
	lauchWorker := make(chan bool, 1)
	workerId := make(chan int, 1)
	childCtx, cancelDeadline := context.WithDeadline(parentCtx, getDeadline())
	defer cancelDeadline()

	lauchWorker <- true

	for newWorker := 1; ; newWorker++ {
		select {
		case <-childCtx.Done():
			fmt.Println("1 RESET WORKER")
			childCtx, cancelDeadline = context.WithDeadline(parentCtx, getDeadline())
			lauchWorker <- true
			break
		case <-lauchWorker:
			fmt.Println("2 START NEW WORKER")
			workerCtx, cancelWorkerCtx := context.WithTimeout(childCtx, getTimeout())
			if newWorker%2 == 0 {
				time.Sleep(1100 * time.Millisecond)
			}

			go worker(workerCtx, workerId)
			result := <-workerId
			cancelWorkerCtx()
			fmt.Println("3 END WORKER RESULT:", result)
			if result == 1 {
				break
			}
		}
	}
}

func worker(ctx context.Context, workerId chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("2.1 cancelled: WID:", workerId)
			workerId <- 1
			return
		default:
			id := 0

			response, err := getCurrenciesFromAPI()
			if err != nil {
				fmt.Println("error:", err.Error())
				id = 1
				workerId <- id
				return
			}

			fmt.Printf("\n\n %+v \n\n", response)

			workerId <- id
			return
		}
	}
}

func getDeadline() time.Time {
	return time.Now().Add(3000 * time.Millisecond)
}

func getTimeout() time.Duration {
	return 1000 * time.Millisecond
}

func getCurrenciesFromAPI() (*service.CurrencyAPIPayload, error) {
	req, err := http.NewRequest("GET", "http://localhost:8000/currencyApiTest", nil)
	if err != nil {
		return &service.CurrencyAPIPayload{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &service.CurrencyAPIPayload{}, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return &service.CurrencyAPIPayload{}, err
	}

	var responseObject service.CurrencyAPIPayload
	json.Unmarshal(bodyBytes, &responseObject)
	return &responseObject, nil
}
