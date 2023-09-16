package main

import (
	"flag"
	"fmt"
)

var noBackLog bool

func main() {
	var number int
	flag.IntVar(&number, "n", 1, "number of numbers to create")

	flag.StringVar(&bucket, "b", bucket, "bucket name")

	flag.BoolVar(&noBackLog, "no-backlog", false, "disable backlog")

	flag.Parse()

	init_db()
	loadcache()
	startSyncing()

	for i := 1; i <= number; i++ {
		n := NewNumber(i)
		fmt.Println("n:", n)
	}

	close(backlog)
	wg.Wait()
}
