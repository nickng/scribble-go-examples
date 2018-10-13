package main

import (
	"flag"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci"
	"github.com/nickng/scribble-go-examples/17_fibonacci/internal/fibonacci"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	i := flag.Int("self", 3, "Specify fib ID")
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
	if 2 < *i && *i < K-1 {
		fibonacci.Fi(protocol, K, *i,
			conn, conn.Port(*i*2-2), // Fib[i-2]->
			conn, conn.Port(*i*2-1), // Fib[i-1]->
			conn, "localhost", conn.Port(*i*2+1), // ->Fib[i+1]
			conn, "localhost", conn.Port(*i*2+2), // ->Fib[i+2]
			wg)
	} else if *i == K-1 {
		fibonacci.Fksub1(protocol, K, *i,
			conn, conn.Port(*i*2-2), // Fib[i-2]->
			conn, conn.Port(*i*2-1), // Fib[i-1]->
			conn, "localhost", conn.Port(*i*2+1), // ->Fib[i+1]
			wg)
	}
	wg.Wait()
}
