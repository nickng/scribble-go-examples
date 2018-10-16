package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync"
	"github.com/nickng/scribble-go-examples/13_lesolver/internal/solver"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protoMain := Solver.New()
	protoSync := Sync.New()

	// Port offset:
	// 0..K*K: use flattened ordinal

	wg := new(sync.WaitGroup)
	wg.Add(K * K)
	first, last := session2.XY(1, 1), session2.XY(K, K)
	toIndex := func(x, y int) int {
		return (y-1)*K + x
	}
	for c := first; c.Lte(last); c = c.Inc(last) {
		switch c.X {
		case 1:
			// will be spawned by W1iWKi
		case K:
			go func(c session2.Pair) {
				solver.W1iWKi(protoMain, protoSync, last, c,
					conn, conn.Port(toIndex(c.X, c.Y)), // W[K-1]->
					conn, "localhost", conn.Port(toIndex(1, c.Y)+1), // ->W[row1+(1,0)]
					conn, conn.Port(last.X*last.Y),
					make(chan struct{}),
					wg)
			}(c)
		default:
			go func(c session2.Pair) {
				time.Sleep(time.Duration((K-c.X)*1000) * time.Millisecond)
				solver.Wii(protoMain, last, c,
					conn, conn.Port(toIndex(c.X, c.Y)),
					conn, "localhost", conn.Port(toIndex(c.X, c.Y)+1),
					wg)
			}(c)
		}
	}
	wg.Wait()
}
