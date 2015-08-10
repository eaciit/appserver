package appserver

import (
	"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/errorlib"
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

	RfcObject interface{}

	listener net.Listener
}

func (a *AppServer) Start() error {
	if a.RfcObject == nil {
		return errorlib.Error(packageName, objAppServer, "Start", "RFC Object is not yet properly initialized")
	}
	a.ReadConfig()

	if a.ServerAddress == "" {
		a.ServerAddress = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
	}

	rpc.Register(a.RfcObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s:%d", a.ServerAddress))
	if e != nil {
		return e
	}

	a.listener = l
	return nil
}

func (a *AppServer) Serve() error {
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
}

func (a *AppServer) Stop() error {
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
		a.Port = config.GetDefault("host", 7890).(int)
	}
	return nil
}
