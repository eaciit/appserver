package appserver

import (
	"errors"
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
	UserID    string
	LoginDate time.Time

	client      *rpc.Client
	secret      string
	referenceID string
}

func (a *Client) Connect(address string, secret string, userid string) error {
	a.UserID = userid
	client, e := rpc.Dial("tcp", address)
	if e != nil {
		return errorlib.Error(packageName, objClient, "Connect", "Unable to connect: "+e.Error())
	}
	a.client = client
	a.LoginDate = time.Now().UTC()

	r := a.Call("addsession", toolkit.M{}.Set("auth_secret", secret).Set("auth_referenceid", a.UserID))
	if r.Status != toolkit.Status_OK {
		return errors.New("Connect: " + r.Message)
	}
	m := toolkit.M{}
	toolkit.FromBytes(r.Data.([]byte), "gob", &m)
	a.secret = m.GetString("secret")
	return nil
}

func (a *Client) Close() {
	if a.client != nil {
		a.Call("removesession", nil)
		a.client.Close()
	}
}

func (a *Client) Call(methodName string, in toolkit.M) *toolkit.Result {
	if a.client == nil {
		return toolkit.NewResult().SetErrorTxt("Unable to call, no connection handshake")
	}
	if in == nil {
		in = toolkit.M{}
	}
	out := toolkit.NewResult()
	in["method"] = methodName
	in["auth_referenceid"] = a.UserID
	if in.Has("auth_secret") == false {
		in.Set("auth_secret", a.secret)
	}
	e := a.client.Call("Rpc.Do", in, out)
	//_ = "breakpoint"
	if e != nil {
		return out.SetError(e)
	}
	return out
}
