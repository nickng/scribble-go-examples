package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
)

// HeadReq is a Scribble HTTP HEAD request wrapper.
type HeadReq struct {
	url     string
	Request *http.Request // Request as the server sees it.
}

// Head creates and wraps a new HTTP HEAD request.
func Head(url string) HeadReq {
	return HeadReq{url: url}
}

// GetOp returns a dummy HeadReq label.
func (HeadReq) GetOp() string {
	return "http.HeadReq"
}

// GetReq is a Scribble HTTP GET request wrapper.
type GetReq struct {
	url      string
	from, to int
	Request  *http.Request // Request as the server sees it.
}

// Get creates and wraps a new HTTP GET request.
func Get(url string, from, to int) GetReq {
	return GetReq{url: url, from: from, to: to}
}

// GetOp returns a dummy GetReq label.
func (GetReq) GetOp() string {
	return "http.GetReq"
}

// Response is a Scribble HTTP response wrapper.
// The body of the response is stored in the Body field.
type Response struct {
	Body    []byte
	Header  http.Header
	Request *http.Request // original request.
}

// GetOp returns a dummy Response label.
func (Response) GetOp() string {
	return "http.Response"
}

var debug = os.Getenv("DEBUG") == "1"

// Formatter is a custom formatter for HTTP requests and responses.
type Formatter struct {
	c         transport2.BinChannel
	emptyBody bool
}

// Wrap wraps a binary channel for encoding/decoding.
func (f *Formatter) Wrap(c transport2.BinChannel) {
	f.c = c
}

// Serialize encodes a Message m as a HTTP request.
func (f *Formatter) Serialize(m session2.ScribMessage) error {
	switch m := m.(type) {
	case *HeadReq:
		req, err := http.NewRequest(http.MethodHead, m.url, nil)
		if err != nil {
			return fmt.Errorf("cannot create HEAD request: %v", err)
		}
		if debug {
			buf := new(bytes.Buffer)
			req.Write(buf)
			log.Println("---- start HTTP debug ---")
			fmt.Fprintf(os.Stderr, buf.String())
			log.Println("---- end HTTP debug ---")
		}
		f.emptyBody = true // set HEAD request flag
		return req.Write(f.c)
	case *GetReq:
		req, err := http.NewRequest(http.MethodGet, m.url, nil)
		if err != nil {
			return fmt.Errorf("cannot create GET request: %v", err)
		}
		if m.from < m.to {
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", m.from, m.to))
		}
		if debug {
			buf := new(bytes.Buffer)
			req.Write(buf)
			log.Println("---- start HTTP debug ---")
			fmt.Fprintf(os.Stderr, buf.String())
			log.Println("---- end HTTP debug ---")
		}
		return req.Write(f.c)
	}
	return fmt.Errorf("invalid message: %#v", m)
}

// Deserialize decodes a HTTP response as a Message m.
func (f *Formatter) Deserialize(m *session2.ScribMessage) error {
	res, err := http.ReadResponse(bufio.NewReader(f.c), nil)
	if err != nil {
		return err
	}
	if debug {
		log.Println("---- start HTTP debug ---")
		for k, v := range res.Header {
			fmt.Fprintf(os.Stderr, "  %s = %s\n", k, v)
		}
		log.Println("---- end HTTP debug ---")
	}
	if f.emptyBody {
		*m = &Response{Header: res.Header}
		f.emptyBody = false // unset HEAD request flag
		return nil
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read HTTP response body: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("cannot close HTTP response body: %v", err)
	}
	*m = &Response{Body: b, Header: res.Header}
	return nil
}