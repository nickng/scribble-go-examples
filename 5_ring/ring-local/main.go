//rhu@HZHL4 ~/code/go
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/5_ring/ring-local-win/main.go
//$ go install github.com/nickng/scribble-go-examples/5_ring/ring-local
//$ bin/ring-local.exe -K=4

package main

import "github.com/nickng/scribble-go-examples/5_ring/internal/ring"

func main() {
	ring.Local()
}
