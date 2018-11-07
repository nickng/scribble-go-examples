package main

import (
	"log"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci"
	"github.com/nickng/scribble-go-examples/17_fibonacci/internal/fibonacci"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Fibonacci.New()

	// Port offsets:
	// 4,5         : Fib[1]->Fib[3], Fib[2]->Fib[3]
	// i*2-2,i*2-1 : Fib[i-2]->Fib[i], Fib[i-2]->Fib[i]
	// K*2-2,K*2-1 : Fib[K-2]->Fib[K], Fib[K-1]->Fib[K]

	if K < 4 {
		log.Fatalf("cannot start with K=%d (need K >= 4)", K)
	}

	wg := new(sync.WaitGroup)
	wg.Add(K)
	for i := 1; i <= K; i++ {
		if i == 1 {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				fibonacci.F1(protocol, K, i,
					conn, "localhost", conn.Port(i*2+2), // ->Fib[3] (i+1)
					wg)
			}(i)
		} else if i == 2 {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				fibonacci.F2(protocol, K, i,
					conn, "localhost", conn.Port(i*2+1), // ->Fib[3] (i+1)
					conn, "localhost", conn.Port(i*2+2), // ->Fib[4] (i+2)
					wg)
			}(i)
		} else if 2 < i && i < K-1 {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				fibonacci.Fi(protocol, K, i,
					conn, conn.Port(i*2-2), // Fib[i-2]->
					conn, conn.Port(i*2-1), // Fib[i-1]->
					conn, "localhost", conn.Port(i*2+1), // ->Fib[i+1]
					conn, "localhost", conn.Port(i*2+2), // ->Fib[i+2]
					wg)
			}(i)
		} else if i == K-1 {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				fibonacci.Fksub1(protocol, K, i,
					conn, conn.Port(i*2-2), // Fib[i-2]->
					conn, conn.Port(i*2-1), // Fib[i-1]->
					conn, "localhost", conn.Port(i*2+1), // ->Fib[i+1]
					wg)
			}(i)
		} else {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				fibonacci.Fk(protocol, K, i,
					conn, conn.Port(i*2-2), // Fib[i-2]->
					conn, conn.Port(i*2-1), // Fib[i-1]->
					wg)
			}(i)
		}
	}
	wg.Wait()
}
