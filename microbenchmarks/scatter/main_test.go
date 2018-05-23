package main_test

import (
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/A_1to1"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/B_1toN"
)

const (
	MinN = 1
	MaxN = 7
)

type rolesScatter struct {
	prot *Scatter.Scatter
	A    *A_1to1.A_1to1
	BN   []*B_1toN.B_1toN
}
