// +build windows

//$ go generate C:/cygwin64/home/rhu/code/go/src/github.com/nickng/scribble-go-examples/9_pget/foreach/gen_windows.go
//$ go install github.com/nickng/scribble-go-examples/9_pget/foreach
//$ c:\cygwin64\home\rhu\code\go> bin\foreach.exe -K=2 http://www.open.ou.nl/ssj/popl19/

//go:generate scribblec-param.bat -d ../ ../PGet.scr -param Sync github.com/nickng/scribble-go-examples/9_pget/PGet -param-api A -param-api B
//go:generate scribblec-param.bat -d ../ ../PGet.scr -param Foreach github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S

package main
