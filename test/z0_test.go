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
	server = new(appserver.Server)
	server.Register(new(controller))
	e := server.Start("localhost:8000")
	if e == nil {
		serverInit = true
	}
}

func TestClient(t *testing.T) {
	checkTestSkip(t)
	client := new(appserver.Client)
	e := client.Connect(server.Address)
	if e != nil {
		t.Error(e.Error())
		return
	}
	defer func() {
		client.Close()
	}()

	var result *toolkit.Result
	result = client.Call("ping", toolkit.M{})
	if result.Status != toolkit.Status_OK {
		t.Error(result.Message)
		return

	} else {
		t.Log("Result: ", result.Data)
	}
}

func TestStop(t *testing.T) {
	checkTestSkip(t)
	server.Stop()
}
