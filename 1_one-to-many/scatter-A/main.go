//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-A
//$ bin/scatter-A.exe

package main

import (
	"encoding/gob"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	var data messages.Data
	gob.Register(&data)
}

func main() {
	connAB, K := scributil.ParseFlags()

	p := Scatter.New()
	scatter.Client_scatter(p, K, 1, connAB, "localhost", connAB.BasePort())
}
