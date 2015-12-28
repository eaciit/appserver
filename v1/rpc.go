package appserver

import (
	"github.com/eaciit/toolkit"
	//"time"
	"errors"
	"strings"
)

//type RpcFn func(toolkit.M, *toolkit.Result) error
type RpcFn func(toolkit.M) *toolkit.Result
type RpcFns map[string]RpcFn

//type ReturnedBytes []byte

type Rpc struct {
	Fns               RpcFns
	Server            *Server
	MarshallingMethod string
}

var _marshallingMethod string

func MarshallingMethod() string {
	if _marshallingMethod == "" {
		_marshallingMethod = "json"
	} else {
		_marshallingMethod = strings.ToLower(_marshallingMethod)
	}
	return _marshallingMethod
}

func SetMarshallingMethod(m string) {
	_marshallingMethod = m
}

func (r *Rpc) Do(in toolkit.M, out *[]byte) error {
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}

	in.Set("rpc", r)
	method := in.GetString("method")
	if method == "" {
		return errors.New("Method is empty")
	}
	fn, fnExist := r.Fns[method]
	if !fnExist {
		return errors.New("Method " + method + " is not exist")
	}
	res := fn(in)
	if res.Status != toolkit.Status_OK {
		return errors.New(res.Message)
	}
	*out = toolkit.ToBytes(res.Data, MarshallingMethod())
	return nil
}

func AddFntoRpc(r *Rpc, svr *Server, k string, fn RpcFn) {
	//func (r *Rpc) AddFn(k string, fn RpcFn) {
	if r.Server == nil {
		r.Server = svr
	}
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}
	r.Fns[k] = fn
}
