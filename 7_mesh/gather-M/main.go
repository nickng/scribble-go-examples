package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Gather"
	"github.com/nickng/scribble-go-examples/7_mesh/internal/mesh"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := Gather.New()
	KK := session2.XY(K, K)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	mesh.MGather(protocol, KK, session2.XY(1, 1), conn, "localhost", conn.BasePort(), wg)
	wg.Wait()
}
