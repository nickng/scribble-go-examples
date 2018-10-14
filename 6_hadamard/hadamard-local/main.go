package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard"
	"github.com/nickng/scribble-go-examples/6_hadamard/internal/hadamard"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Hadamard.New()

	max := session2.XY(K, K)

	wg := new(sync.WaitGroup)
	wg.Add(max.Flatten(max) * 3)
	for c := session2.XY(1, 1); c.Lte(max); c = c.Inc(max) {
		go func(c session2.Pair) {
			time.Sleep(100 * time.Millisecond)
			hadamard.A(protocol, max, c, conn, "localhost", conn.Port(c.Flatten(max)), wg)
		}(c)
		go func(c session2.Pair) {
			time.Sleep(100 * time.Millisecond)
			hadamard.B(protocol, max, c, conn, "localhost", conn.Port(max.Flatten(max)+c.Flatten(max)), wg)
		}(c)
		go func(c session2.Pair) {
			hadamard.C(protocol, max, c, conn, conn.Port(c.Flatten(max)), conn, conn.Port(max.Flatten(max)+c.Flatten(max)), wg)
		}(c)
	}
	wg.Wait()
}
