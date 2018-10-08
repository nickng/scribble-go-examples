//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-B
//$ bin/foreach-B.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Foreach github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api B -param-api A

package main

import (
	"encoding/gob"
	"sync"

	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

func main() {
	listen, _, fmtr, port, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	p := Foreach.New()  // FIXME: K should be param here?
	for i := 1; i <= K; i++ {
		go foreach.Server_gather(listen, fmtr, port+i, p, K, i, wg)
	}

	wg.Wait()
}
