package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody"
	"github.com/nickng/scribble-go-examples/12_nbody/internal/nbody"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := NBody.New()

	// Port offsets:
	// 1    : W[1] → W[K]
	// 2..K : W[K-1] → W[K]

	wg := new(sync.WaitGroup)
	wg.Add(K)
	for i := 1; i <= K; i++ {
		if i == 1 {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				nbody.W1(protocol, K, i, conn, "localhost", conn.Port(i+1), conn, "localhost", conn.Port(1), wg)
			}(i)
		} else if 1 < i && i < K {
			go func(i int) {
				time.Sleep(time.Duration((K-i)*100) * time.Millisecond)
				nbody.Wi(protocol, K, i, conn, conn.Port(i), conn, "localhost", conn.Port(i+1), wg)
			}(i)
		} else {
			go nbody.WK(protocol, K, i, conn, conn.Port(i), conn, conn.Port(1), wg)
		}
	}
	wg.Wait()
}
