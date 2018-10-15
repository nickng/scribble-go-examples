package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Diagonal"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

// Shared memory diagonal pipeline

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Diagonal.New()

	// Port offsets
	// 2..K: W[(K-1,K-1)] -> W[(K,K)]

	wg := new(sync.WaitGroup)
	wg.Add(K)
	first, increment, last := session2.XY(1, 1), session2.XY(1, 1), session2.XY(K, K)
	for c := first; c.Lte(last); c = c.Plus(increment) {
		if c.Eq(first) {
			go func(c session2.Pair) {
				time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
				mesh.W11Diag(protocol, last, c,
					conn, "localhost", conn.Port(c.Plus(increment).Flatten(last)), // ->W[(K+1,1)]
					wg)
			}(c)
		} else if first.Lt(c) && c.Lt(last) {
			go func(c session2.Pair) {
				time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
				mesh.WiiDiag(protocol, last, c,
					conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
					conn, "localhost", conn.Port(c.Plus(increment).Flatten(last)), // ->W[(K+1,1)]
					wg)
			}(c)
		} else {
			go func(c session2.Pair) {
				mesh.WKKDiag(protocol, last, c,
					conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
					wg)
			}(c)
		}
	}
	wg.Wait()
}
