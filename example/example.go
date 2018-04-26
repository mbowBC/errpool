package main

import (
	"fmt"
	"time"

	"github.com/stuarthicks/errpool"
)

func main() {
	g := errpool.Group{}
	g.StartWorkers(2)

	for i := 1; i <= 10; i++ {
		i := i
		g.Run(func() error {
			time.Sleep(1 * time.Second)
			return fmt.Errorf("%d failed", i)
		})
	}

	errors := g.Wait()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}
