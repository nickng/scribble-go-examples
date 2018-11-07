package main

import (
	"flag"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi"
	"github.com/nickng/scribble-go-examples/18_jacobi/internal/jacobi"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	self := flag.Int("I", 3, "Specify worker ID")
	conn, K := scributil.ParseFlags()
	protocol := Jacobi.New()

	// Port offsets:
	// 1..K     : W[i-1] -> W[i]
	// K+4..K+K : W[2] -> W[4..K]

	if K < 3 {
		log.Fatalf("cannot start with K=%d (need K >= 3)", K)
	}
	if 2 >= *self || *self >= K {
		log.Fatalf("self (%d) should be 2 < self < K", *self)
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	jacobi.Wi(protocol, K, *self,
		conn, conn.Port(*self), // ->W[k]
		conn, "localhost", conn.Port(*self+1), // W[k]->
		conn, conn.Port(K+*self), // scatter channel
		wg)
	wg.Wait()
}
