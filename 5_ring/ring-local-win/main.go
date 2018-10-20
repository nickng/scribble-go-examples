// +build windows

//C:\cygwin64\home\rhu\code\go>
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/5_ring/ring-local-win/main.go
//$ go install github.com/nickng/scribble-go-examples/5_ring/ring-local
//$ bin\ring-local-win.exe -K=4

//go:generate scribblec-param.bat ../Ring.scr -d ../ -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

package main

import "github.com/nickng/scribble-go-examples/5_ring/internal/ring"

func main() {
	ring.Local()
}
