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
	/*
		times := make([]time.Time, 0)
		for i := 0; i < 1000; i++ {
			times = append(times, time.Now().UTC())
		}
		result.Data = toolkit.GetEncodeByte(times)
		return nil
	*/
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

func (r *Rpc) AddFn(svr *AppServer, k string, fn RpcFn) {
	//func (r *Rpc) AddFn(k string, fn RpcFn) {
	if r.Server != nil {
		r.Server = svr
	}
	if r.Fns == nil {
		r.Fns = map[string]RpcFn{}
	}

	r.Fns[k] = fn
}
