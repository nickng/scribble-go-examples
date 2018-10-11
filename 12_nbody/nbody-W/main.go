package main

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody"
	"github.com/nickng/scribble-go-examples/12_nbody/internal/nbody"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	self := flag.Int("self", 2, "W[i] index (1 < i < K)")
	conn, K := scributil.ParseFlags()
	protocol := NBody.New()

	// Port offsets:
	// 1    : W[1] → W[K]
	// 2..K : W[K-1] → W[K]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	nbody.Wi(protocol, K, *self, conn, conn.Port(*self), conn, "localhost", conn.Port(*self+1), wg)
	wg.Wait()
}
