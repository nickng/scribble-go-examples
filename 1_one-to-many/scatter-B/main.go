//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-B
//$ bin/scatter-B.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B


package main

import (
	"encoding/gob"
	"sync"

	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

func main() {
	listen, _, fmtr, port, K := scributil.ParseFlags()

	p := Scatter.New()  // FIXME: K should be param here?
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 1; i <= K; i++ {
		go scatter.Server_gather(listen, fmtr, port+i, p, K, i, wg)
	}
	wg.Wait()
}
