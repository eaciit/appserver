package appserver

import (
	"errors"
	"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net"
	"net/rpc"
	"reflect"
)

const (
	packageName = "eaciit"
	objServer   = "Server"
)

type IServer interface {
	Start() error
	Stop() error

	ReadConfig() error
}

type Server struct {
	ServerId   string
	ConfigFile string
	ServerName string
	Port       int
	Address    string
	Role       string

	rpcObject interface{}

	Log      *toolkit.LogEngine
	listener net.Listener

	container interface{}
}

func (a *Server) SetContainer(o interface{}) {
	a.container = o
}

func (a *Server) Container() interface{} {
	return a.container
}

/*
func (a *Server) Start(reloadConfig bool) error {
	if a.rpcObject == nil {
		return errorlib.Error(packageName, objServer, "Start", "RPC Object is not yet properly initialized")
	}
	if reloadConfig {
		a.ReadConfig()
	}

	if a.Address == "" {
		a.Address = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
	}

	rpc.Register(a.rpcObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s", a.Address))
	if e != nil {
		return e
	}

	a.listener = l
	return nil
}
*/
func (a *Server) Start(address string) error {
	if a.rpcObject == nil {
		return errorlib.Error(packageName, objServer, "Start", "RPC Object is not yet properly initialized")
	}
	/*
		if reloadConfig {
			a.ReadConfig()
		}
	*/

	if a.Address == "" {
		if address != "" {
			a.Address = address
		} else {
			a.Address = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
		}
		if a.Address == "" {
			return errors.New("RPC Server address is empty")
		}
	}

	if a.Log == nil {
		le, e := toolkit.NewLog(true, false, "", "", "")
		if e == nil {
			a.Log = le
		} else {
			return errors.New("Unable to setup log")
		}
	}
	rpc.Register(a.rpcObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s", a.Address))
	if e != nil {
		return e
	}
	a.listener = l
	go func() {
		rpc.Accept(l)
	}()
	return nil
}

func (a *Server) AddFn(methodname string, fn RpcFn) {
	var r *Rpc
	if a.rpcObject == nil {
		r = new(Rpc)
	} else {
		r = a.rpcObject.(*Rpc)
	}

	AddFntoRpc(r, a, methodname, fn)
	//r.AddFn(a, methodname, f)
	//r.AddFn(methodname, f)
	a.rpcObject = r
}

func (a *Server) Register(o interface{}) error {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		return errors.New("Invalid object for RPC Register")
	}
	methodCount := t.NumMethod()
	for i := 0; i < methodCount; i++ {
		method := t.Method(i)
		mtype := method.Type

		//-- now check method signature
		if mtype.NumIn() == 2 && mtype.In(1).String() == "toolkit.M" {
			if mtype.NumOut() == 1 && mtype.Out(0).String() == "*toolkit.Result" {
				a.AddFn(method.Name, nil)
			}
		}
	}
	return nil
}

func (a *Server) Serve() error {
	/*
		for {
			conn, e := a.listener.Accept()
			if e != nil {
				return e
			}
			go func(c net.Conn) {
				defer c.Close()
				rpc.ServeConn(c)
			}(conn)
		}
	*/
	rpc.Accept(a.listener)
	return nil
}

func (a *Server) Stop() error {
	a.listener.Close()
	a.Log.Info("Stopping service")
	return nil
}

func (a *Server) ReadConfig() error {
	if a.ConfigFile == "" {
		a.ServerName = "localhost"
		a.Port = 7890
	} else {
		e := config.SetConfigFile(a.ConfigFile)
		if e != nil {
			return e
		}

		a.ServerName = config.GetDefault("host", "localhost").(string)
		a.Port = config.GetDefault("port", 7890).(int)
	}
	return nil
}
