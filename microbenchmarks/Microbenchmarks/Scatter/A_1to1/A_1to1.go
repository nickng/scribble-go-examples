package A_1to1
//[A{[1]}, B{[1,N]}]


import "github.com/rhu1/scribble-go-runtime/runtime/session2"
import "github.com/rhu1/scribble-go-runtime/runtime/transport2"

type A_1to1 struct {
Proto session2.Protocol
Self int
*session2.LinearResource
lin uint64
MPChan *session2.MPChan
N int
Params map[string]int
_Init *Init
_End *End
}

func New(p session2.Protocol, N int, self int) *A_1to1 {
ep := &A_1to1{
p,
self,
&session2.LinearResource{},
1,
session2.NewMPChan(self, []string{"A", "B"}),
N,
make(map[string]int),
nil,
nil,
}
ep._Init = &Init{ nil, 1, ep }
ep._End = &End{ nil, 0, ep }
return ep
}

func (ini *A_1to1) B_1toN_Accept(id int, ss transport2.ScribListener, sfmt session2.ScribMessageFormatter) error {
defer ini.MPChan.ConnWg.Done()
ini.MPChan.ConnWg.Add(1)
c, err := ss.Accept()
if err != nil {
return err
}

sfmt.Wrap(c)
ini.MPChan.Conns["B"][id] = c
ini.MPChan.Fmts["B"][id] = sfmt
return err
}

func (ini *A_1to1) B_1toN_Dial(id int, host string, port int, dialler func (string, int) (transport2.BinChannel, error), sfmt session2.ScribMessageFormatter) error {
defer ini.MPChan.ConnWg.Done()
ini.MPChan.ConnWg.Add(1)
c, err := dialler(host, port)
if err != nil {
return err
}

sfmt.Wrap(c)
ini.MPChan.Conns["B"][id] = c
ini.MPChan.Fmts["B"][id] = sfmt
return err
}

func (ini *A_1to1) Run(f func(*Init) End) End {
defer ini.Close()
return f(ini.Init())
}

func (ini *A_1to1) Init() *Init {
ini.Use()
ini.MPChan.CheckConnection()
return ini._Init
}

func (ini *A_1to1) Close() {
defer ini.MPChan.Close()
}