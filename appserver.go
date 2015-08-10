package appserver

import (
	"fmt"
	"github.com/eaciit/config"
	"net"
	"net/rpc"
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

	RcvrObject interface{}
}

func (a *AppServer) Start() (net.Listener, error) {
	if a.ServerAddress == "" {
		a.ServerAddress = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
	}

	rpc.Register(a.RcvrObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s:%d", a.ServerAddress))
	if e != nil {
		return nil, e
	}

	return l, nil
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

		a.ServerName = config.GetDefault("host", "localhost")
		a.Port = config.GetDefault("host", 7890)
	}
}
