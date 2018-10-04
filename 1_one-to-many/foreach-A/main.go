//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-A
//$ bin/foreach-A.exe

package main

import (
	"encoding/gob"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

func main() {
	connAB, K := scributil.ParseFlags()

	p := Foreach.New()
	foreach.Client_scatter(p, K, 1, connAB, "localhost", connAB.BasePort())
}
