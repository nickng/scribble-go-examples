package A_1toM

import "github.com/rhu1/scribble-go-runtime/runtime/session2"
import "github.com/rhu1/scribble-go-runtime/runtime/transport2"

type A_1toM struct {
Proto session2.Protocol
Self int
*session2.LinearResource
lin uint64
MPChan *session2.MPChan
N int
M int
Params map[string]int
_Init *Init
_End *End
}

func New(p session2.Protocol, N int, M int, self int) *A_1toM {
ep := &A_1toM{
p,
self,
&session2.LinearResource{},
1,
session2.NewMPChan(self, []string{"A", "B"}),
N,
M,
make(map[string]int),
nil,
nil,
}
ep._Init = &Init{ nil, 1, ep }
ep._End = &End{ nil, 0, ep }
return ep
}

func (ini *A_1toM) B_1toN_Accept(id int, ss transport2.ScribListener, sfmt session2.ScribMessageFormatter) error {
defer ini.MPChan.ConnWg.Done()
ini.MPChan.ConnWg.Add(1)
c, err := ss.Accept()
sfmt.Wrap(c)
ini.MPChan.Conns["B"][id] = c
ini.MPChan.Fmts["B"][id] = sfmt
return err
}

func (ini *A_1toM) B_1toN_Dial(id int, host string, port int, dialler func (string, int) (transport2.BinChannel, error), sfmt session2.ScribMessageFormatter) error {
defer ini.MPChan.ConnWg.Done()
ini.MPChan.ConnWg.Add(1)
c, err := dialler(host, port)
sfmt.Wrap(c)
ini.MPChan.Conns["B"][id] = c
ini.MPChan.Fmts["B"][id] = sfmt
return err
}

func (ini *A_1toM) Run(f func(*Init) End) End {
//defer ini.MPChan.Close()
ini.Use()
ini.MPChan.CheckConnection()
return f(ini._Init)
}
