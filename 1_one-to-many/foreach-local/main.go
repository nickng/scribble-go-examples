//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-local
//$ bin/foreach-local.exe -t=shm

package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(K)

	p := Foreach.New() // FIXME: K should be param here?
	for i := 1; i <= K; i++ {
		go foreach.Server_gather(p, K, i, connAB, connAB.Port(i), wg)
	}
	time.Sleep(100 * time.Millisecond)
	foreach.Client_scatter(p, K, 1, connAB, "localhost", connAB.BasePort())
	wg.Wait()
}
