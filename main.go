package main

import (
	"currency-go/app"

	_ "github.com/lib/pq"
)

const APIKEY = "heLJSwAL4RUBRYfPaxmTVXvmz4tp2t3YWip8mITQ"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "currency"
)

func main() {
	app.Run()
	// // connection string
	// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// // open database
	// db, err := sql.Open("postgres", psqlconn)
	// if err != nil {
	// 	panic(err)
	// }
	// // close database
	// defer db.Close()

	// // check db
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Connected!")
}
