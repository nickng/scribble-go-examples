//rhu@HZHL4 ~/code/go
//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local/main.go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local
//$ bin/scatter-local.exe -t=shm -K=4

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package main

import "github.com/nickng/scribble-go-examples/1_one-to-many/scatter"

func main() {
	scatter.Local()
}
