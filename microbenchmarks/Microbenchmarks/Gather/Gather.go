package Gather

import A_1toN "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/A_1toN"
import B_1to1 "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/B_1to1"

type Gather struct {
}

func (*Gather) IsProtocol() {
}

func New() *Gather {
return &Gather{ }
}

func (p *Gather) New_A_1toN(N int, self int) *A_1toN.A_1toN {
return A_1toN.New(p, N, self)
}

func (p *Gather) New_B_1to1(N int, self int) *B_1to1.B_1to1 {
return B_1to1.New(p, N, self)
}
