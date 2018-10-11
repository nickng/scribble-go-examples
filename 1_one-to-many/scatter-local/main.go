//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local
//$ bin/scatter-local.exe -t=shm

package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K, _ := scributil.ParseFlags()

	p := Scatter.New() // FIXME: K should be param here?
	wg := new(sync.WaitGroup)
	wg.Add(K)
	for i := 1; i <= K; i++ {
		go scatter.Server_gather(p, K, i, connAB, connAB.Port(i), wg)
	}
	time.Sleep(100 * time.Millisecond)
	scatter.Client_scatter(p, K, 1, connAB, "localhost", connAB.BasePort())
	wg.Wait()
}
