//rhu@HZHL4 ~/code/go
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/5_ring/ring-local/main.go
//$ go install github.com/nickng/scribble-go-examples/5_ring/ring-local
//$ bin/ring-local.exe -K=4

//go:generate scribblec-param.sh ../Ring.scr -d ../ -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

package main

import "github.com/nickng/scribble-go-examples/5_ring/internal/ring"

func main() {
	ring.Local()
}
