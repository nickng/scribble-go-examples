//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local
//$ bin/scatter-local.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B


package main

import (
	"encoding/gob"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

//..HERE: take CL arg for transport
func main() {
	listen, dial, fmtr, port, K := scributil.ParseFlags()

	p := Scatter.New()  // FIXME: K should be param here?
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 1; i <= K; i++ {
		go scatter.Server_gather(listen, fmtr, port+i, p, K, i, wg)
	}
	time.Sleep(100 * time.Millisecond)
	scatter.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
	wg.Wait()
}
