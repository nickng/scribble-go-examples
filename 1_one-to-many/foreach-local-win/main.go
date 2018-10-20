// +build windows

//rhu@HZHL4 ~/code/go
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/1_one-to-many/foreach-local-win/main.go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/foreach-local
//$ bin\foreach-local.exe -t=shm -K=4

//go:generate scribblec-param.bat ../OneToMany.scr -d ../ -param Foreach github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import (
	"github.com/nickng/scribble-go-examples/1_one-to-many/foreach"
)

func main() {
	foreach.Local()
}
