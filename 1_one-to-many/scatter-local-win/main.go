// +build windows

//C:\cygwin64\home\rhu\code\go>
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local-win/main.go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local-win
//$ bin\scatter-local-win.exe -t=shm -K=4

//go:generate scribblec-param.bat ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import (
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

func main() {
	scatter.Local()
}