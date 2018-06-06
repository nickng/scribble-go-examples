package main_test

import (
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/A_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/B_1to1"
)

const (
	MinN = 1
	MaxN = 7
)

type rolesGather struct {
	prot *Gather.Gather
	AN   []*A_1toN.A_1toN
	B    *B_1to1.B_1to1
}
