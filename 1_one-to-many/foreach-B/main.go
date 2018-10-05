//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-B
//$ bin/foreach-B.exe

package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	p := Foreach.New() // FIXME: K should be param here?
	for i := 1; i <= K; i++ {
		go foreach.Server_gather(p, K, i, connAB, connAB.Port(i), wg)
	}
	wg.Wait()
}
