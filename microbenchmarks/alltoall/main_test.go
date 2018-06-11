//go:generate scribblec-param.sh ../Microbenchmarks.scr -d ../ -param Alltoall github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks -param-api A -param-api B

package main

import (
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/A_1toM"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/B_1toN"
)

const (
	MinM = 1
	MaxM = 6
	MinN = 1
	MaxN = 6
)

type rolesAlltoall struct {
	prot *Alltoall.Alltoall
	AM   []*A_1toM.A_1toM
	BN   []*B_1toN.B_1toN
}
