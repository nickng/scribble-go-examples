//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-A
//$ bin/foreach-A.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Foreach github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import (
	"encoding/gob"

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
	_, dial, fmtr, port, K := scributil.ParseFlags()

	p := Foreach.New()
	foreach.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
}
