package main

import (
	"context"
	"fmt"
	"time"
)

func test() {
	doSomething()
}

func doSomething() {
	numberCh := make(chan int)
	isBreak := false
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancelCtx()

	go doAnother(ctx, numberCh)

	for num := 1; ; num++ {
		select {
		case numberCh <- num:
			go doAnother(ctx, numberCh)
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			isBreak = true
			break
		}
		if isBreak {
			break
		}

	}

	cancelCtx()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("doSomething: finished\n")
}

func doAnother(ctx context.Context, printInt chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("doAnother err: %s\n", err)
			}
			fmt.Printf("doAnother: finished\n")
			return
		case num := <-printInt:
			fmt.Printf("doAnother: %d\n", num)
		}
	}
}
