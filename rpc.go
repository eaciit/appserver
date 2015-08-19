package appserver

import (
	"github.com/eaciit/toolkit"
	//"time"
)

type RpcFn func(toolkit.M, *toolkit.Result) error
type RpcFns map[string]RpcFn

type Rpc struct {
	Fns    RpcFns
	Server *AppServer
}

func (r *Rpc) Do(in toolkit.M, result *toolkit.Result) error {
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}

	if in.Has("method") {
		in.Set("rpc", r)
		fn := r.Fns[in.Get("method", "").(string)]
		return fn(in, result)
	}
	return nil
}

func AddFntoRpc(r *Rpc, svr *AppServer, k string, fn RpcFn) {
	//func (r *Rpc) AddFn(k string, fn RpcFn) {
	if r.Server == nil {
		r.Server = svr
	}
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}
	r.Fns[k] = fn
}
