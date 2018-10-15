package main

import (
	"flag"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver"
	"github.com/nickng/scribble-go-examples/13_lesolver/internal/solver"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	x := flag.Int("x", 2, "x-coordinate")
	y := flag.Int("y", 2, "y-coordinate")
	conn, K := scributil.ParseFlags()
	protoMain := Solver.New()

	// Port offset:
	// 0..K*K: use flattened ordinal

	wg := new(sync.WaitGroup)
	wg.Add(1)
	last := session2.XY(K, K)
	toIndex := func(x, y int) int {
		return (y-1)*K + x
	}
	func(c session2.Pair) {
		time.Sleep(time.Duration((K-c.X)*100) * time.Millisecond)
		solver.Wii(protoMain, last, c,
			conn, conn.Port(toIndex(c.X, c.Y)),
			conn, "localhost", conn.Port(toIndex(c.X, c.Y)+1),
			wg)
	}(session2.XY(*x, *y))
	wg.Wait()
}
