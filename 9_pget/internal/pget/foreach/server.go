package foreach

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach"
	S "github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach/family_1/S_1to1"
	httpmsg "github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/nickng/scribble-go-examples/9_pget/internal/pget"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

// RunServer runs the server at URL accepting K threads.
func RunServer(K int, URL string) {
	_, httpPort := pget.ParseURL(URL)

	protocol := Foreach.New()
	S := protocol.New_family_1_S_1to1(K, 1)

	ln, err := tcp.Listen(httpPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()

	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(K)
	mu := new(sync.Mutex)
	for i := 1; i <= K; i++ {
		if i == 1 {
			go func(i int) {
				scributil.Debugf("S: listening for F[%d] at :%d.\n", i, httpPort)
				mu.Lock()
				if err := S.F_1to1and1toK_Accept(i, ln, new(HTTPFormatter)); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
				mu.Unlock()
				wgSvr.Done()
			}(i)
		} else {
			go func(i int) {
				scributil.Debugf("S: listening for F[%d] at :%d.\n", i, httpPort)
				mu.Lock()
				if err := S.F_1toK_not_1to1_Accept(i, ln, new(HTTPFormatter)); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
				mu.Unlock()
				wgSvr.Done()
			}(i)
		}
	}
	wgSvr.Wait()
	scributil.Debugf("S: Ready.\n")
	S.Run(serverBody)
}

func serverBody(s *S.Init) S.End {
	req := make([]httpmsg.HeadReq, 1)
	s0 := s.F_1_Gather_Head(req)
	scributil.Debugf("S: received %v.\n", req)
	res := []httpmsg.Response{
		httpmsg.Response{
			Body:    []byte(content()),
			Header:  map[string][]string{http.CanonicalHeaderKey("Content-Length"): []string{strconv.Itoa(len(content()))}},
			Request: req[0].Request, // Request being handled.
		},
	}
	s1 := s0.F_1_Scatter_Res(res)
	scributil.Debugf("S: sent %v.\n", res)
	sEnd := s1.Foreach(nestedS)
	return *sEnd
}

func nestedS(s *S.Init_49) S.End {
	req := make([]httpmsg.GetReq, 1)
	s0 := *s.F_I_Gather_Get(req)
	res := []httpmsg.Response{
		httpmsg.Response{
			Body:    []byte(content()),
			Header:  map[string][]string{http.CanonicalHeaderKey("Content-Length"): []string{strconv.Itoa(len(content()))}},
			Request: req[0].Request, // Request being handled.
		},
	}
	sEnd := s0.F_I_Scatter_Res(res)
	return *sEnd
}

// HTTPFormatter is a server-side HTTP formatter.
type HTTPFormatter struct {
	c transport2.BinChannel
}

// Wrap wraps a server-side TCP connection.
func (f *HTTPFormatter) Wrap(c transport2.BinChannel) {
	f.c = c
}

// Serialize emulates sending of a file requested.
func (f *HTTPFormatter) Serialize(m session2.ScribMessage) error {
	resp := m.(*httpmsg.Response)
	var from, to int
	if resp.Request.Header.Get("Range") != "" {
		re := regexp.MustCompile(`bytes=(\d+)-(\d+)?`)
		matches := re.FindStringSubmatch(resp.Request.Header.Get("Range"))
		if len(matches) == 3 {
			from, _ = strconv.Atoi(matches[1])
			to, _ = strconv.Atoi(matches[2])
		} else if len(matches) == 2 {
			from, _ = strconv.Atoi(matches[1])
		}
	}
	var (
		body []byte
		blen int
	)
	if resp.Request.Method == http.MethodHead {
		// HEAD requests return the full length
		// but the body is empty.
		body, blen = []byte{}, len(resp.Body)
	} else {
		// GET requests returns the partial length
		// with corresponding partition of the body.
		if to < from {
			body = resp.Body[from:]
			blen = len(body)
		} else {
			body = resp.Body[from : to+1]
			blen = len(body)
		}
	}
	res := &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       resp.Request,
		ContentLength: int64(blen),
		Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
		Header: http.Header{
			"Content-Length": []string{strconv.Itoa(blen)},
			"Content-Type":   []string{"text/html"},
		},
	}

	scributil.Debugf("--- start server HTTP debug ---\n")
	var buf bytes.Buffer
	res.Write(&buf)
	scributil.Debugf("%s\n", buf.String())
	scributil.Debugf("--- end server HTTP debug ---\n")
	_, err := buf.WriteTo(f.c)
	return err
}

// Deserialize emulates server reading an HTTP Request.
func (f *HTTPFormatter) Deserialize(m *session2.ScribMessage) error {
	req, err := http.ReadRequest(bufio.NewReader(f.c))
	if err != nil {
		return err
	}
	switch req.Method {
	case http.MethodHead:
		*m = &httpmsg.HeadReq{Request: req} // Head
		return nil
	case http.MethodGet:
		*m = &httpmsg.GetReq{Request: req} // Res
		return nil
	}
	return fmt.Errorf("method %s not handled", req.Method)
}

func content() string {
	return `<html><!DOCTYPE html><body><a href='http://www.open.ou.nl/ssj/popl19/'>Scribble-Go</a></body></html>`
}
