package main

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync"
	"github.com/nickng/scribble-go-examples/13_lesolver/internal/solver"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	y := flag.Int("y", 1, "y-coordinate")
	conn, K := scributil.ParseFlags()
	protoMain := Solver.New()
	protoSync := Sync.New()

	// Port offset:
	// 0..K*K: use flattened ordinal

	wg := new(sync.WaitGroup)
	wg.Add(1)
	last := session2.XY(K, K)
	toIndex := func(x, y int) int {
		return (y-1)*K + x
	}
	func(c session2.Pair) {
		solver.W1i(protoMain, protoSync, last, c,
			conn, "localhost", conn.Port(toIndex(c.X, c.Y)+1),
			conn, conn.Port(last.X*last.Y),
			wg)
	}(session2.XY(1, *y))
	wg.Wait()
}
