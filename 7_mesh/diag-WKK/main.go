package main

import (
	"sync"

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
	wg.Add(1)
	last := session2.XY(K, K)
	func(c session2.Pair) {
		mesh.WKKDiag(protocol, last, c,
			conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
			wg)
	}(last)
	wg.Wait()
}
