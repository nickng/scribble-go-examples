package http_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	httpmsg "github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
)

type customBuf struct {
	bytes.Buffer
}

func (b *customBuf) Close() error { return nil }

// This tests that the formatter can send standard HTTP messages.
func TestFormatterSend(t *testing.T) {
	var c transport2.BinChannel = new(customBuf)
	f := new(httpmsg.Formatter)
	f.Wrap(c)

	const requrl = "http://example.com/path"
	matchURL := func(u *url.URL, req *http.Request) bool {
		return u.Host == req.Host && u.Path == req.RequestURI
	}

	t.Run("HEAD", func(t *testing.T) {
		if err := f.Serialize(httpmsg.Head(requrl)); err != nil {
			t.Fatalf("serialise failed: %v", err)
		}
		req, err := http.ReadRequest(bufio.NewReader(c))
		if err != nil {
			t.Fatal(err)
		}
		u, err := url.Parse(requrl)
		if err != nil {
			t.Fatal(err)
		}
		if req.Method != http.MethodHead || !matchURL(u, req) {
			t.Fatalf("mismatched request: HEAD %s -- %v", u.String(), req)
		}
	})
	t.Run("GET", func(t *testing.T) {
		if err := f.Serialize(httpmsg.Get(requrl, 0, 255)); err != nil {
			t.Fatalf("serialise failed: %v", err)
		}
		req, err := http.ReadRequest(bufio.NewReader(c))
		if err != nil {
			t.Fatal(err)
		}
		u, err := url.Parse(requrl)
		if err != nil {
			t.Fatal(err)
		}
		if req.Method != http.MethodGet || !matchURL(u, req) {
			t.Fatalf("mismatched request: GET %s -- %v", u.String(), req)
		}
	})
}

// This tests that the formatter can receive standard HTTP messages.
func TestFormatterRecv(t *testing.T) {
	var c transport2.BinChannel = new(customBuf)
	f := new(httpmsg.Formatter)
	f.Wrap(c)

	const body = "[response body]"

	res := &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Body:          ioutil.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}
	res.Write(c)

	var msg session2.ScribMessage = new(httpmsg.Response)
	if err := f.Deserialize(&msg); err != nil {
		t.Fatalf("deserialise failed: %v", err)
	}
	rcvdBytes := msg.(*httpmsg.Response).Body
	t.Logf("received: %v", rcvdBytes)
	if want, got := body, string(rcvdBytes); want != got {
		t.Fatalf("mismatched response body: %s", got)
	}
}
