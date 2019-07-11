package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	_ "net/http/pprof"

	leak "github.com/zimmski/go-leak"
)

var x map[string]int

func BenchmarkEmptyMapCap100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = make(map[string]int, 100)
	}
}

func ExampleCache() {
	dir, _ := ioutil.TempDir("", "example")
	fmt.Println(dir)
	return

	cache.CacheInit()
	created := time.Now()

	for i := 0; i < 10000; i++ {
		segment := new(Segment)
		segment.BasePage = 16
		segment.Count = 128

		segment.Created = time.Now()
		segment.Transaction = nil
		segment.TxResult = nil

		b := make([]byte, 10)
		rand.Read(b)

		segments,_ := cache.Get(string(b))
		segments.Push(segment)
		cache.Set(string(b), segments)
	}

	fmt.Println("Add Tx to Cache")
	fmt.Printf("Cache: %+v\n\n", cache)

	fmt.Println("Append to Segment")
	fmt.Println("Run Life")

	ticker := time.NewTicker(time.Second)
	for {
		t := <-ticker.C
		if t.Sub(created).Seconds() > 20.0 {
			ticker.Stop()
			break
		} else {
			fmt.Printf("Cache: %+v\n\n", cache)
		}
	}

	// Output:
}

var iter = 10000000
var visible = false

func k_main() {
	cache.CacheInit()
	fmt.Printf("Iter -> : %d\n", iter)

	for i := 0; i < 10; i++ {
		fmt.Printf("-------------------[Case-%02d]-------------------\n", i+1)
		ww := time.Now()
		fmt.Println(ww)
		q := leak.MarkMemory()
		// test()
		leaksq := q.Release()

		fmt.Printf("[%d] Leak Result: %d\n", i+1, leaksq)
		fmt.Println("Time Elepsed:", time.Since(ww))
		visible = false
		// fmt.Println("-----------------------------------------------\n\n")
	}
}

// func k_test() {
// 	//defer profile.Start(profile.MemProfile).Stop()
// 	//created := time.Now()

// 	q := leak.MarkMemory()
// 	for i := 0; i < iter; i++ {
// 		q := make(SegmentMap, 100)
// 		q[0] = Segment{}
// 	}
// 	leaksq := q.Release()

// 	fmt.Printf("[Map] Leak -> : %d\n", leaksq)
// 	visible = true
// 	m := leak.MarkMemory()

// 	for i := 0; i < iter; i++ {
// 		s := cache.Get(fmt.Sprintf("test-%d", i))
// 		s.Push(&Segment{}, cache.ticker)
// 		cache.Set(fmt.Sprintf("test-%d", i), s)
// 	}
// 	time.Sleep(20 * time.Second)
// 	leaks := m.Release()

// 	// for i := 0; i < 3000000; i++ {
// 	// 	segment := new(Segment)
// 	// 	segment.BasePage = 16
// 	// 	segment.Count = 128

// 	// 	segment.Created = time.Now()

// 	// 	b := make([]byte, 10)
// 	// 	rand.Read(b)

// 	// 	segments := cache.Get(string(b))
// 	// 	segments.Push(segment, cache.ticker)
// 	// 	cache.Set(string(b), segments)
// 	// }

// 	// fmt.Println("Add Tx to Cache")
// 	// //fmt.Printf("Cache: %+v\n\n", cache)

// 	// fmt.Println("Append to Segment")
// 	// fmt.Println("Run Life")

// 	// ticker := time.NewTicker(time.Second)
// 	// for {
// 	// 	t := <-ticker.C
// 	// 	if t.Sub(created).Seconds() > 60.0 && len(cache.Mat) == 0 {
// 	// 		ticker.Stop()
// 	// 		//fmt.Printf("Cache: %+v\n\n", cache)
// 	// 		break
// 	// 	} else {
// 	// 		fmt.Printf("data: %d\n", len(cache.Mat))
// 	// 	}
// 	// }

// 	fmt.Printf("[Cache] Leak -> : %d\n", leaks)
// 	//runtime.GC()

// }
