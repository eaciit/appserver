package appserver

import (
	"github.com/eaciit/toolkit"
	//"time"
	"errors"
)

type RpcFn func(toolkit.M, *toolkit.Result) error
type RpcFns map[string]RpcFn

type Rpc struct {
	Fns    RpcFns
	Server *Server
}

func (r *Rpc) Do(method string, in toolkit.M, result *toolkit.Result) error {
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}

	in.Set("rpc", r)
	fn, fnExist := r.Fns[method]
	if !fnExist {
		return errors.New("Method " + method + " is not exist")
	}
	return fn(in, result)
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
