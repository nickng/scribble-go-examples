package main

import (
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi"
	"github.com/nickng/scribble-go-examples/18_jacobi/internal/jacobi"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Jacobi.New()

	// Port offsets:
	// 1..K     : W[i-1] -> W[i]
	// K+4..K+K : W[2] -> W[4..K]

	if K < 3 {
		log.Fatalf("cannot start with K=%d (need K >= 3)", K)
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	jacobi.W2(protocol, K, 2,
		conn, conn.Port(2), // ->W[k]
		conn, "localhost", conn.Port(3), // W[k]->
		conn, "localhost", conn.Port(K), // scatter channel
		wg)
	wg.Wait()
}
