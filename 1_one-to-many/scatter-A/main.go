//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-A
//$ bin/scatter-A.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import (
	"encoding/gob"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	//"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

//var _ = shm.Dial
var _ = tcp.Dial

func init() {
	var data onetomany.Data
	gob.Register(&data)
}

func main() {
	port := 33333
	K := 2
	p := Scatter.New()
	dial := tcp.Dial
	fmtr := func() session2.ScribMessageFormatter { return new(session2.GobFormatter) } 
	/*// Not applicable to distributed scenario
	dial := shm.Dial
	fmtr := func() session2.ScribMessageFormatter { return new(session2.PassByPointer) }*/
	scatter.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
}
