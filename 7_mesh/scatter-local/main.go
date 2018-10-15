package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Scatter"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Scatter.New()
	KK := session2.XY(K, K)

	wg := new(sync.WaitGroup)
	wg.Add(KK.Flatten(KK) + 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		mesh.MScatter(protocol, KK, session2.XY(1, 1), conn, "localhost", conn.BasePort(), wg)
	}()
	for w := session2.XY(1, 1); w.Lte(KK); w = w.Inc(KK) {
		go func(w session2.Pair) {
			mesh.WScatter(protocol, KK, w, conn, conn.Port(w.Flatten(KK)), wg)
		}(w)
	}
	wg.Wait()
}
