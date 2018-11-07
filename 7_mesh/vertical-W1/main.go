package main

// Shared memory column pipeline

import (
	"sync"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Mesh3.New()

	// Port offsets
	// 2..K: W[(K-1,1)] -> W[(K,1)]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	first, last := session2.XY(1, 1), session2.XY(1, K)
	func(c session2.Pair) {
		mesh.W11v(protocol, last, c,
			conn, "localhost", conn.Port(c.Flatten(last)+1), // ->W[(K+1,1)]
			wg)
	}(first)
	wg.Wait()
}
