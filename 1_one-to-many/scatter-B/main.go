//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-B
//$ bin/scatter-B.exe

package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(K)

	p := Scatter.New() // FIXME: K should be param here?
	for i := 1; i <= K; i++ {
		go scatter.Server_gather(p, K, i, connAB, connAB.Port(i), wg)
	}
	wg.Wait()
}
