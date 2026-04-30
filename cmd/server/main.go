package main

import (
	"fmt"
	"time"
)

func randomFunc() string {
	return "This is the function ill return inside the goroutine\n"
}

func main() {
	fmt.Printf("THis is the main entrypoint of the server\n")
	go func() {
		res := randomFunc()
		fmt.Printf("%v", res)
	}()

	time.Sleep(time.Second * 1)

}
