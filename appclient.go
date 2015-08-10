package appserver

import (
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net/rpc"
	"time"
)

const (
	objAppClient = "AppClient"
)

type AppClient struct {
	UserId    string
	FullName  string
	SessionId string
	LoginDate time.Time

	client *rpc.Client
}

func (a *AppClient) Connect(address string) error {
	client, e := rpc.Dial("tcp", address)
	if e != nil {
		return errorlib.Error(packageName, objAppClient, "Connect", "Unable to connect: "+e.Error())
	}
	a.client = client
	a.LoginDate = time.Now().UTC()
	return nil
}

func (a *AppClient) Close() {
	if a.client != nil {
		a.client.Close()
	}
}

func (a *AppClient) Call(methodName string, in toolkit.M) *toolkit.Result {
	out := new(toolkit.Result)
	out.Status = toolkit.Status_OK
	start := time.Now()
	e := a.client.Call("Rpc."+methodName, in, out)
	out.Duration = time.Since(start)
	if e != nil {
		out.Status = toolkit.Status_NOK
		out.Message = e.Error()
	}
	return out
}
