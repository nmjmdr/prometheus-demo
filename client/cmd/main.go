package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	done := make(chan bool)
	go makeCalls(done)
	fmt.Println("Press enter to exit")
	fmt.Scanln()
	done <- true
}

func makeCalls(done chan bool) {
	g, _ := errgroup.WithContext(context.Background())
	maxParallel := 200
	minParallel := 50
	for {
		fmt.Println("Running again")
		for i := 0; i < (rand.Intn(maxParallel-minParallel) + minParallel); i++ {
			g.Go(func() error {
				rs, err := http.Get("http://localhost:1323/ping")
				if err != nil {
					return err
				}
				defer rs.Body.Close()
				b, err := ioutil.ReadAll(rs.Body)
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			fmt.Println("Error: ", err)
		}

		select {
		case <-done:
			return
		default:
			continue
		}
	}
}
