//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-local
//$ bin/foreach-local.exe -t=shm

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B


package main

import (
	"encoding/gob"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

// FIXME: -t=shm crashes
func main() {
	listen, dial, fmtr, port, K := scributil.ParseFlags()

	p := Foreach.New()  // FIXME: K should be param here?
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 1; i <= K; i++ {
		go foreach.Server_gather(listen, fmtr, port+i, p, K, i, wg)
	}
	time.Sleep(100 * time.Millisecond)
	foreach.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
	wg.Wait()
}
