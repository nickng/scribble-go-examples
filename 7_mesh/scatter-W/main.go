package main

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Scatter"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	x := flag.Int("x", 2, "x-coordinate of worker")
	y := flag.Int("y", 2, "y-coordinate of worker")
	conn, K := scributil.ParseFlags()
	protocol := Scatter.New()
	KK := session2.XY(K, K)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	func(w session2.Pair) {
		mesh.WScatter(protocol, KK, w, conn, conn.Port(w.Flatten(KK)), wg)
	}(session2.XY(*x, *y))
	wg.Wait()
}
