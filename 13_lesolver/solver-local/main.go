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
			go func(c session2.Pair) {
				time.Sleep(time.Duration((K-c.X)*100) * time.Millisecond)
				solver.W1i(protoMain, protoSync, last, c,
					conn, "localhost", conn.Port(toIndex(c.X, c.Y)+1),
					conn, conn.Port(last.X*last.Y),
					wg)
			}(c)
		case K:
			go func(c session2.Pair) {
				time.Sleep(time.Duration((K-c.X)*100) * time.Millisecond)
				solver.WKi(protoMain, last, c,
					conn, conn.Port(toIndex(c.X, c.Y)),
					wg)
			}(c)
		default:
			go func(c session2.Pair) {
				time.Sleep(time.Duration((K-c.X)*100) * time.Millisecond)
				solver.Wii(protoMain, last, c,
					conn, conn.Port(toIndex(c.X, c.Y)),
					conn, "localhost", conn.Port(toIndex(c.X, c.Y)+1),
					wg)
			}(c)
		}
	}
	wg.Wait()
}
