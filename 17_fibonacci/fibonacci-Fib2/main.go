package main

import (
	"log"
	"sync"

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
	wg.Add(1)
	fibonacci.F2(protocol, K, 2,
		conn, "localhost", conn.Port(5), // ->Fib[3] (i+1)
		conn, "localhost", conn.Port(6), // ->Fib[4] (i+2)
		wg)
	wg.Wait()
}
