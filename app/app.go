package app

import (
	"currency-go/domain"
	"currency-go/service"
	"fmt"
	"net/http"

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
	client := NewClientDB()
	defer client.Close()

	ch := NewCurrencyHandler(service.NewDefaultCurrencyService(domain.NewCurrencyRepositoryStub(client)))

	router := mux.NewRouter()
	router.HandleFunc("/currencies/{id}", ch.Get).Methods(http.MethodGet)
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
