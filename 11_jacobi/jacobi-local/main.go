package main

import (
	"log"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/11_jacobi/Jacobi/Jacobi"
	"github.com/nickng/scribble-go-examples/11_jacobi/internal/jacobi"
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
	wg.Add(K)
	for k := 1; k <= K; k++ {
		if k == 1 {
			go func(k int) {
				time.Sleep(time.Duration((K-k)*100) * time.Millisecond)
				jacobi.W1(protocol, K, k, conn, "localhost", conn.Port(k+1), wg)
			}(k)
		} else if k == 2 {
			go func(k int) {
				time.Sleep(time.Duration((K-k)*100) * time.Millisecond)
				jacobi.W2(protocol, K, k,
					conn, conn.Port(k), // ->W[k]
					conn, "localhost", conn.Port(k+1), // W[k]->
					conn, "localhost", conn.Port(K), // scatter channel
					wg)
			}(k)
		} else if 2 < k && k < K {
			go func(k int) {
				time.Sleep(time.Duration((K-k)*100) * time.Millisecond)
				jacobi.Wi(protocol, K, k,
					conn, conn.Port(k), // ->W[k]
					conn, "localhost", conn.Port(k+1), // W[k]->
					conn, conn.Port(K+k), // scatter channel
					wg)
			}(k)
		} else {
			go func(k int) {
				time.Sleep(time.Duration((K-k)*100) * time.Millisecond)
				jacobi.Wk(protocol, K, k,
					conn, conn.Port(k), // ->W[k]
					conn, conn.Port(K+k), // scatter channel
					wg)
			}(k)
		}
	}
	wg.Wait()
}
