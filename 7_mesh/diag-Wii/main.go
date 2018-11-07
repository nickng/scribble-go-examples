package main

import (
	"flag"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Diagonal"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

// Shared memory diagonal pipeline

func main() {
	xy := flag.Int("xy", 2, "Specify xy-coordinate of worker, e.g. -xy 2 = (2,2)")
	conn, K := scributil.ParseFlags()
	protocol := Diagonal.New()

	// Port offsets
	// 2..K: W[(K-1,K-1)] -> W[(K,K)]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	increment, last := session2.XY(1, 1), session2.XY(K, K)
	func(c session2.Pair) {
		time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
		mesh.WiiDiag(protocol, last, c,
			conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
			conn, "localhost", conn.Port(c.Plus(increment).Flatten(last)), // ->W[(K+1,1)]
			wg)
	}(session2.XY(*xy, *xy))
	wg.Wait()
}
