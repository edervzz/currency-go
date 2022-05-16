package main

import (
	"context"
	"currency-go/app"
	"fmt"
	"time"

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

func do() {
	parentCtx := context.Background()
	lauchWorker := make(chan bool, 1)
	workerId := make(chan int, 1)
	newWorker := 0
	childCtx, cancelChildCtx := context.WithTimeout(parentCtx, 4500*time.Millisecond)
	defer cancelChildCtx()

	lauchWorker <- true

	for ; ; newWorker++ {

		select {
		case <-childCtx.Done():
			newWorker = 0
			fmt.Println("RESET WORKERS:", newWorker)
			childCtx, cancelChildCtx = context.WithTimeout(parentCtx, 4000*time.Millisecond)
			lauchWorker <- true
			break
		case <-lauchWorker:
			fmt.Println("START NEW WORKER", newWorker)
			go worker(childCtx, workerId)
			result := <-workerId
			fmt.Println("END WORKER RESULT:", result)
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
			fmt.Println("cancelled: WID:", workerId)
			workerId <- 1
			return
		default:
			id := 0
			fmt.Printf("Processing...: %d\n", id)
			time.Sleep(900 * time.Millisecond)
			workerId <- id
			return
		}
	}
}
