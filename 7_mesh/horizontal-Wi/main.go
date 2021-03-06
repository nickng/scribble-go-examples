package main

// shared memory horizontal wave.

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	x := flag.Int("x", 2, "Specify x-coordinate")
	conn, K := scributil.ParseFlags()
	protocol := Mesh1.New()

	// Port offsets
	// 2..K: W[(K-1,1)] -> W[(K,1)]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	first, last := session2.XY(1, 1), session2.XY(K, 1)
	func(c session2.Pair) {
		mesh.Wi1(protocol, last, c,
			conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
			conn, "localhost", conn.Port(c.Flatten(last)+1), // ->W[(K+1,1)]
			wg)
	}(first.Plus(session2.XY(*x-1, 0)))
	wg.Wait()
}
