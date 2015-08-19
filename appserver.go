package appserver

import (
	"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net"
	"net/rpc"
)

const (
	packageName  = "eaciit"
	objAppServer = "AppServer"
)

type IServer interface {
	Start() error
	Stop() error

	ReadConfig() error
}

type AppServer struct {
	ServerId      string
	ConfigFile    string
	ServerName    string
	Port          int
	ServerAddress string
	Role          string

	RpcObject interface{}

	Log      *toolkit.LogEngine
	listener net.Listener
}

func (a *AppServer) Start(reloadConfig bool) error {
	if a.RpcObject == nil {
		return errorlib.Error(packageName, objAppServer, "Start", "RPC Object is not yet properly initialized")
	}
	if reloadConfig {
		a.ReadConfig()
	}

	if a.ServerAddress == "" {
		a.ServerAddress = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
	}

	rpc.Register(a.RpcObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s", a.ServerAddress))
	if e != nil {
		return e
	}

	a.listener = l
	return nil
}

func (a *AppServer) AddFn(methodname string, f RpcFn) {
	var r *Rpc
	if a.RpcObject == nil {
		r = new(Rpc)
	} else {
		r = a.RpcObject.(*Rpc)
	}
	r.AddFn(a, methodname, f)
	//r.AddFn(methodname, f)
	a.RpcObject = r
}

func (a *AppServer) Serve() error {
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

func (a *AppServer) Stop() error {
	a.Log.Info("Stopping service")
	return nil
}

func (a *AppServer) ReadConfig() error {
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
