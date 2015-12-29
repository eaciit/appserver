package rpctest

import (
	"github.com/eaciit/appserver/v1"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

var server *appserver.Server
var client *appserver.Client
var serverInit bool
var (
	serverSecret string = "ariefdarmawan"
)

type controller struct {
}

type Score struct {
	Subject string
	Value   int
}

func (a *controller) Hi(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	name := in.GetString("Name")
	r.SetBytes(struct {
		HelloMessage string
		TimeNow      time.Time
		Scores       []Score
	}{"Hello " + name, time.Now(), []Score{{"Bahasa Indonesia", 90}, {"Math", 85}}}, "gob")
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
	server.SetSecret(serverSecret)
	e := server.Start("localhost:8000")
	if e == nil {
		serverInit = true
	}
}

func checkResult(result *toolkit.Result, t *testing.T) {
	if result.Status != toolkit.Status_OK {
		t.Error(result.Message)
	} else {
		if result.IsEncoded() == false {
			t.Logf("Result: %v", result.Data)
		} else {

			m := struct {
				HelloMessage string
				TimeNow      time.Time
				Scores       []Score
			}{}

			//m := toolkit.M{}
			e := result.GetFromBytes(&m)
			if e != nil {
				t.Errorf("Unable to decode result: %s\n", e.Error())
				return
			}
			t.Logf("Result (decoded): %s", toolkit.JsonString(m))
		}
	}
}

func TestClient(t *testing.T) {
	checkTestSkip(t)
	client = new(appserver.Client)
	e := client.Connect(server.Address, serverSecret, "ariefdarmawan_10")
	if e != nil {
		t.Error(e.Error())
		return
	}

	var result *toolkit.Result
	result = client.Call("ping", toolkit.M{})
	checkResult(result, t)
}

func TestClientHi(t *testing.T) {
	checkTestSkip(t)
	r := client.Call("hi", toolkit.M{})
	checkResult(r, t)
}

func TestStop(t *testing.T) {
	checkTestSkip(t)
	server.Stop()
	if client != nil {
		client.Close()
	}
}
