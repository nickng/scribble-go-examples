package main

import (
	"sync"

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
	wg.Add(1)
	nbody.W1(protocol, K, 1, conn, "localhost", conn.Port(2), conn, "localhost", conn.Port(1), wg)
	wg.Wait()
}
