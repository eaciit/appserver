package appserver

import (
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net/rpc"
	"time"
)

const (
	objClient = "Client"
)

var _timeOut time.Duration

func SetDialTimeout(t time.Duration) {
	_timeOut = t
}

func DialTimeout() time.Duration {
	if _timeOut == 0 {
		_timeOut = 90 * time.Second
	}
	return _timeOut
}

type Client struct {
	UserId    string
	FullName  string
	SessionId string
	LoginDate time.Time

	client *rpc.Client
}

func (a *Client) Connect(address string) error {
	client, e := rpc.Dial("tcp", address)
	if e != nil {
		return errorlib.Error(packageName, objClient, "Connect", "Unable to connect: "+e.Error())
	}
	a.client = client
	a.LoginDate = time.Now().UTC()
	return nil
}

func (a *Client) Close() {
	if a.client != nil {
		a.client.Close()
	}
}

func (a *Client) Call(methodName string, in toolkit.M) *toolkit.Result {
	out := toolkit.NewResult()
	in["method"] = methodName
	e := a.client.Call("Rpc.Do", in, out)
	//_ = "breakpoint"
	if e != nil {
		return out.SetError(e)
	}
	return out
}
