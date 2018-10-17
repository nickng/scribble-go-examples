// Package pget is a internal helper package for pget implementations.
// It contains common helper functions to both versions of pget.
package pget

import (
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/nickng/scribble-go-examples/9_pget/msg"
)

// ParseURL parses a URL and splits to host and port part
// for use in establishing connection.
func ParseURL(URL string) (host string, port int) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("invalid URL: %v", err)
	}
	if host = u.Hostname(); host == "" {
		log.Fatalf("invalid host: %s", URL)
	}
	if u.Scheme == "https" {
		log.Fatalf("invalid URL: https:// URL not supported")
	}
	if port, err = strconv.Atoi(u.Port()); err != nil {
		port, err = net.LookupPort("tcp", u.Scheme)
		if err != nil {
			log.Fatal(err)
		}
	}
	return host, port
}

var debug = os.Getenv("DEBUG") == "1"

// Debugf is a debug print function.
func Debugf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}

// AllocResponse allocates memory to store received Response.
func AllocResponse() []http.Response {
	return make([]http.Response, 1)
}

// MakeHeadReq returns a new HTTP HEAD request.
func MakeHeadReq(url string) []http.HeadReq {
	return []http.HeadReq{http.Head(url)}
}

// ExtractSize extracts the content size from a HEAD response.
func ExtractSize(res []http.Response) int {
	size, err := strconv.Atoi(res[0].Header.Get("Content-Length"))
	if err != nil {
		return 0
	}
	return size
}

// MakeMeta returns a new Meta message for Master
// given the url and size of the content.
func MakeMeta(url string, size int) []msg.Meta {
	return []msg.Meta{msg.Meta{URL: url, Size: size}}
}

// AllocMeta allocates memory to store received Meta message.
func AllocMeta() []msg.Meta {
	return make([]msg.Meta, 1)
}

// MakeJobs constructs fetcher jobs from a given Meta message meta
// and the given number of fetcher K.
func MakeJobs(meta []msg.Meta, K int) []msg.Job {
	jobs := make([]msg.Job, K)
	fragSize := meta[0].Size / K
	for i := 0; i < K; i++ {
		jobs[i].URL = meta[0].URL
		jobs[i].RangeFrom = i * fragSize
		if i < K-1 {
			jobs[i].RangeTo = (i+1)*fragSize - 1
		} else {
			// The last fetcher fetches the rest.
			jobs[i].RangeTo = meta[0].Size
		}
	}
	return jobs
}

// AllocJob allocates memory to store received Job message.
func AllocJob() []msg.Job {
	return make([]msg.Job, 1)
}

// MakeGetReq returns a new HTTP GET request.
func MakeGetReq(job []msg.Job) []http.GetReq {
	return []http.GetReq{http.Get(job[0].URL, job[0].RangeFrom, job[0].RangeTo)}
}

// CollectData collates the HTTP response received and
// returns a Data message for the Master.
func CollectData(res []http.Response) []msg.Data {
	return []msg.Data{msg.Data{Data: res[0].Body}}
}

// AllocData allocates memory to store the received Data.
func AllocData(K int) []msg.Data {
	return make([]msg.Data, K)
}

// MakeDone creates a Done signal.
func MakeDone() []msg.Done {
	return []msg.Done{msg.Done{}}
}

// AllocDone allocates memory to store the received Done signal.
func AllocDone() []msg.Done {
	return make([]msg.Done, 1)
}
