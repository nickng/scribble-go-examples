package http

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
)

// HeadReq is a Scribble HTTP HEAD request wrapper.
type HeadReq struct {
	url string
}

// Head creates and wraps a new HTTP HEAD request.
func Head(url string) HeadReq {
	return HeadReq{url}
}

// GetOp returns a dummy HeadReq label.
func (HeadReq) GetOp() string {
	return "http.HeadReq"
}

// GetReq is a Scribble HTTP GET request wrapper.
type GetReq struct {
	url      string
	from, to int
}

// Get creates and wraps a new HTTP GET request.
func Get(url string, from, to int) GetReq {
	return GetReq{url, from, to}
}

// GetOp returns a dummy GetReq label.
func (GetReq) GetOp() string {
	return "http.GetReq"
}

// Response is a Scribble HTTP response wrapper.
// The body of the response is stored in the Body field.
type Response struct {
	Body   []byte
	Header http.Header
}

// GetOp returns a dummy Response label.
func (Response) GetOp() string {
	return "http.Response"
}

// Formatter is a custom formatter for HTTP requests and responses.
type Formatter struct {
	c transport2.BinChannel
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
			return fmt.Errorf("cannot create GET request: %v", err)
		}
		return req.Write(f.c)
	case *GetReq:
		req, err := http.NewRequest(http.MethodGet, m.url, nil)
		if err != nil {
			return fmt.Errorf("cannot create GET request: %v", err)
		}
		if m.from < m.to {
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", m.from, m.to))
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read HTTP response: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("cannot close HTTP response body: %v", err)
	}
	*m = Response{Body: b, Header: res.Header}
	return nil
}
