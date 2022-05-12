package main

import (
	"log"

	"github.com/jcuvillier/workerpool"
)

func main() {
	wp := workerpool.New[int](5, false)
	for i := 0; i < 100; i++ {
		wp.Add(func() (int, error) {
			return 1, nil
		})
	}
	k := 0
	for r := range wp.Exec() {
		k += r
	}
	if err := wp.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println(k) // print 100
}
