// +build windows

//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/9_pget/server/gen_windows.go
//$ go install github.com/nickng/scribble-go-examples/9_pget/server

//c:\cygwin64\home\rhu\code\go>
//$ bin\server.exe -K=2 http://www.dummy.com

//go:generate scribblec-param.bat -d ../ ../PGet.scr -param Foreach github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S -parforeach

package main
