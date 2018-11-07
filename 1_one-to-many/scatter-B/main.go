//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-B
//$ bin/scatter-B.exe

package main

import (
	"sync"
	"flag"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
        var I int
	flag.IntVar(&I, "I", -1, "self ID (1 <= I <= K)")
	connAB, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(1)

	p := Scatter.New() // FIXME: K should be param here?
	go scatter.Server_gather(p, K, I, connAB, connAB.Port(I), wg)
	wg.Wait()
}
