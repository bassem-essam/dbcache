package main

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type Number struct {
	gorm.Model
	Value int

	// A session-like field to group numbers together
	Bucket string
}

var NumberCache = sync.Map{}
var NextIDCache = sync.Map{}

func (n Number) String() string {
	return fmt.Sprint(n.Value)
}

func NewNumber(i int) *Number {
	if n, ok := NumberCache.Load(i); ok {
		fmt.Printf("NumberCache.Load(%d)\n", i)
		return n.(*Number)
	}

	n := &Number{Value: i, Bucket: bucket}
	NumberCache.Store(i, n)

	fmt.Printf("NumberCache.Store(%d)\n", i)

	return CreateNumber(n)
}

var numbersNextID chan uint = make(chan uint)

// GetNextID returns the next available ID for the given table.
func GetNextID(table string) uint {
	if id, ok := NextIDCache.Load(table); ok {
		NextIDCache.Store(table, id.(uint)+1)
		return id.(uint)
	}

	var n Number
	res := db.Model(&Number{}).Last(&n)
	if res.Error == gorm.ErrRecordNotFound {
		n.ID = 0
	} else if res.Error != nil {
		panic(res.Error)
	}

	// Store the second next ID in the cache to be ready for the next call
	NextIDCache.Store(table, n.ID+2)

	return n.ID + 1
}

// This way of insertion is too slow when used for a large number of objects.
func SlowInsert(n *Number) *Number {
	res := db.Create(n)
	if res.Error != nil {
		panic(res.Error)
	}

	return n
}

func CreateNumber(n *Number) *Number {
	if noBackLog {
		return SlowInsert(n)
	}

	// Set the ID to the next available ID
	// This is a hack to mimic the behavior of the database insertion
	// To allow you to use the object as if it was inserted, now it can be linked to other
	// tables with foreign keys for example.
	// The real insertion will be done in the background by the goroutine (the backlog)
	backlog <- n
	n.ID = <-numbersNextID
	// fmt.Println("numbersNextID <-", n.ID)
	return n
}

var backlog = make(chan *Number, 100)
var wg sync.WaitGroup

// startSyncing starts a goroutine that will synchronize the backlog with db
// i.e. the NumberCache with the database.
func startSyncing() {
	wg.Add(1)

	var numbers []*Number
	go func() {
		for n := range backlog {
			numbersNextID <- GetNextID("numbers")

			numbers = append(numbers, n)

			if len(numbers) > 1000 {
				res := db.Model(&Number{}).Create(numbers)
				if res.Error != nil {
					panic(res.Error)
				}
				numbers = nil
			}
		}

		if len(numbers) > 0 {
			res := db.Model(&Number{}).Create(numbers)
			if res.Error != nil {
				panic(res.Error)
			}
		}

		wg.Done()
	}()
}
