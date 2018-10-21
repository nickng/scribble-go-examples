//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/9_pget/foreach/win/main.go
//$ go install github.com/nickng/scribble-go-examples/9_pget/foreach/win
//$ c:\cygwin64\home\rhu\code\go>bin\win.exe -K=2 http://www.example.com

//go:generate scribblec-param.bat -d ../../ ../../PGet.scr -param Sync github.com/nickng/scribble-go-examples/9_pget/PGet -param-api A -param-api B
//go:generate scribblec-param.bat -d ../../ ../../PGet.scr -param Foreach github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S

package main

import "flag"
import "fmt"
import "log"
import "os"

import "github.com/nickng/scribble-go-examples/9_pget/internal/pget/foreach"


var (
	K int
	URL string
)

func init() {
	flag.IntVar(&K, "K", 2, "Specify number of fetchers")
	log.SetPrefix("pget: ")
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pget [-K fetchers] URL\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	URL = flag.Arg(0)

	foreach.RunClient(K, URL)
}
