package main

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard"
	"github.com/nickng/scribble-go-examples/6_hadamard/internal/hadamard"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	x := flag.Int("x", 1, "x-coordinate")
	y := flag.Int("y", 1, "y-coordinate")
	conn, K := scributil.ParseFlags()
	protocol := Hadamard.New()

	max := session2.XY(K, K)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	func(c session2.Pair) {
		hadamard.C(protocol, max, c, conn, conn.Port(c.Flatten(max)), conn, conn.Port(max.Flatten(max)+c.Flatten(max)), wg)
	}(session2.XY(*x, *y))
	wg.Wait()
}
