package main

// shared memory horizontal wave.

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Mesh1.New()

	// Port offsets
	// 2..K: W[(K-1,1)] -> W[(K,1)]

	wg := new(sync.WaitGroup)
	wg.Add(K)
	first, increment, last := session2.XY(1, 1), session2.XY(1, 0), session2.XY(K, 1)
	for c := first; c.Lte(last); c = c.Plus(increment) {
		if c.Eq(first) {
			go func(c session2.Pair) {
				time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
				mesh.W11(protocol, last, c,
					conn, "localhost", conn.Port(c.Flatten(last)+1), // ->W[(K+1,1)]
					wg)
			}(c)
		} else if first.Lt(c) && c.Lt(last) {
			go func(c session2.Pair) {
				time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
				mesh.Wi1(protocol, last, c,
					conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
					conn, "localhost", conn.Port(c.Flatten(last)+1), // ->W[(K+1,1)]
					wg)
			}(c)
		} else {
			go func(c session2.Pair) {
				mesh.WK1(protocol, last, c,
					conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
					wg)
			}(c)
		}
	}
	wg.Wait()
}
