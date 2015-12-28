package rpctest

import (
	"github.com/eaciit/appserver/v1"
	"github.com/eaciit/toolkit"
	"testing"
)

var server *appserver.Server
var appclient *appserver.Client
var serverInit bool

type controller struct {
}

func (a *controller) Hi(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	name := in.GetString("Name")
	r.Data = "Hello " + name
	return r
}

func checkTestSkip(t *testing.T) {
	if serverInit == false {
		t.Skip()
	}
}

func TestStart(t *testing.T) {
	defer func() {
		serverInit = true
	}()
	server = new(appserver.Server)
	server.Register(new(controller))
	server.Start("localhost:8800")
}

func TestStop(t *testing.T) {
	checkTestSkip(t)
}
