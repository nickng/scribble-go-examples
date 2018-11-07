package main

// Shared memory column pipeline

import (
	"flag"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	y := flag.Int("y", 2, "Specify y-coordinate")
	conn, K := scributil.ParseFlags()
	protocol := Mesh3.New()

	// Port offsets
	// 2..K: W[(K-1,1)] -> W[(K,1)]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	first, last := session2.XY(1, 1), session2.XY(1, K)
	func(c session2.Pair) {
		time.Sleep(time.Duration((last.Flatten(last)-c.Flatten(last))*100) * time.Millisecond)
		mesh.W1iv(protocol, last, c,
			conn, conn.Port(c.Flatten(last)), // W[(K-1,1)]->
			conn, "localhost", conn.Port(c.Flatten(last)+1), // ->W[(K+1,1)]
			wg)
	}(first.Plus(session2.XY(0, *y-1)))
	wg.Wait()
}
