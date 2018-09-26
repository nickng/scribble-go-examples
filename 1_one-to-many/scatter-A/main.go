//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-A
//$ bin/scatter-A.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import (
	"encoding/gob"

	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

func init() {
	var data onetomany.Data
	gob.Register(&data)
}

func main() {
	_, dial, fmtr, port, K := scributil.ParseFlags()

	p := Scatter.New()
	scatter.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
}
