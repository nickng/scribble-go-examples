package Alltoall

import A_1toM "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/A_1toM"
import B_1toN "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/B_1toN"

type Alltoall struct {
}

func (*Alltoall) IsProtocol() {
}

func New() *Alltoall {
return &Alltoall{ }
}

func (p *Alltoall) New_A_1toM(M int, N int, self int) *A_1toM.A_1toM {
return A_1toM.New(p, M, N, self)
}

func (p *Alltoall) New_B_1toN(M int, N int, self int) *B_1toN.B_1toN {
return B_1toN.New(p, M, N, self)
}
