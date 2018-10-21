package foreach

import (
	"log"

	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach"
	S "github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach/family_1/S_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/nickng/scribble-go-examples/9_pget/internal/pget"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	//"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

var (
	// URL is the URL to fetch.
	URLserver string
)

func RunServer(K int, URL string) {
	URLserver = URL

	httpHost, httpPort := pget.ParseURL(URL)

	protocol := Foreach.New()
	S := protocol.New_family_1_S_1to1(K, 1)

	ln, err := tcp.Listen(80)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()

	for i := 1; i <= K; i++ {
		if i == 1 {
			if err := S.F_1to1and1toK_Accept(i, ln, new(HTTPFormatter)); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		} else {
			if err := S.F_1toK_not_1to1_Accept(i, ln, new(HTTPFormatter)); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		}
	}

	S.Run(serverBody)
}

func serverBody(s *S.Init) S.End {
	s0 := s.F_1_Gather_Head(..TODO..)
	s1 := s0.F_1_Scatter_Res(..TODO..)
	/*s2 := s1.F_1toK_Gather_Get()
	sEnd := s2.F_1toK_Scatter_Res()*/
	sEnd := s1.Foreach(nestedS)
	return *sEnd
}

func nestedS(s *S.Init_49) S.End {
	return *s.F_I_Gather_Get(..TODO..).F_I_Scatter_Res(..TODO..)
}

// HTTPFormatter is a server-side HTTP formatter.
type HTTPFormatter struct {
	c transport2.BinChannel
}

// Wrap wraps a server-side TCP connection.
func (f *HTTPFormatter) Wrap(c transport2.BinChannel) { f.c = c }

// Serialize emulates sending of a file requested.
func (f *HTTPFormatter) Serialize(m session2.ScribMessage) error {
	file := `Content of HTTP file`
	res := &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          ioutil.NopCloser(strings.NewReader(file)),
		ContentLength: int64(len(file)),
	}
	return res.Write(f.c)
}

// Deserialize emulates server reading an HTTP Request.
func (f *HTTPFormatter) Deserialize(m *session2.ScribMessage) error {
	_, err := http.ReadRequest(bufio.NewReader(f.c))
	if err != nil {
		return err
	}
	return nil
}
