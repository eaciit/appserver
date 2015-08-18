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

func (a *AppClient) Call(methodName string, in toolkit.M, result interface{}) error {
	out := new(toolkit.Result)
	in["method"]=methodName
	e := a.client.Call("Rpc.Do", in, out)
	//_ = "breakpoint"
	if e != nil {
		return errorlib.Error(packageName, objAppClient, "Call", e.Error())
	} else if out.Status == toolkit.Status_NOK {
		return errorlib.Error(packageName, objAppClient, "Call", out.Message)
	} else {
		if result != nil {
			if out.Data!=nil {
				result = toolkit.DecodeByte(out.Data.([]byte), result)
			}
		}
	}
	return nil
}
