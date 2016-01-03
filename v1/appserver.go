package appserver

import (
	"errors"
	"fmt"
	//"github.com/eaciit/config"
	//"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net"
	"net/rpc"
	"reflect"
	"strings"
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
	//ServerName string
	//Port       int

	Address           string
	Role              string
	UseGlobalPassword bool
	AllowMultiLogin   bool

	rpcObject *Rpc

	Log      *toolkit.LogEngine
	listener net.Listener

	container interface{}
	users     map[string]*User
	sessions  map[string]*Session
	secret    string
}

func (a *Server) SetContainer(o interface{}) {
	a.container = o
}

func (a *Server) Container() interface{} {
	return a.container
}

func (a *Server) SetSecret(s string) {
	a.secret = s
}

func (a *Server) Secret() string {
	if a.secret == "" {
		a.secret = toolkit.GenerateRandomString("", 32)
	}
	return a.secret
}

func (a *Server) validateSecret(secretType string, referenceID string, secret string) bool {
	secretType = strings.ToLower(secretType)
	//referenceID = referenceID
	if secretType == "" || secretType == "self" {
		user, userExist := a.users[referenceID]
		if userExist == false {
			return false
		} else if a.UseGlobalPassword == true {
			return secret == a.Secret()
		} else {
			return secret == user.Secret
		}
	} else if secretType == "session" {
		//a.Log.Info(fmt.Sprintf("Session Validation: ID=%s Secret=%s\n%s", referenceID, secret, toolkit.JsonString(a.sessions)))
		session, exist := a.sessions[referenceID]
		if !exist {
			a.Log.Warning("Session " + referenceID + " could not be found on " + a.Address)
			return false
		}
		//if session.ReferenceID != referenceID {
		//	return false
		//}
		if session.IsValid() == false {
			//a.Log.Warning("Session "+referenceID+" is not valid")
			return false
		}
		return session.Secret == secret
	}
	return false
}

func (a *Server) Start(address string) error {
	if a.Log == nil {
		le, e := toolkit.NewLog(true, false, "", "", "")
		if e == nil {
			a.Log = le
		} else {
			return errors.New("Unable to setup log")
		}
	}

	if a.rpcObject == nil {
		//return errorlib.Error(packageName, objServer, "Start", "RPC Object is not yet properly initialized")
		a.rpcObject = new(Rpc)
	}
	/*
		if reloadConfig {
			a.ReadConfig()
		}
	*/

	if a.Address == "" {
		if address != "" {
			a.Address = address
		}
		/*else {
			a.Address = fmt.Sprintf("%s:%d", a.ServerName, a.Port)
		}
		*/
		if a.Address == "" {
			return errors.New("RPC Server address is empty")
		}
	}

	rpc.Register(a.rpcObject)
	l, e := net.Listen("tcp", fmt.Sprintf("%s", a.Address))
	if e != nil {
		return e
	}

	//init a ping method. Ping method will return "EACIIT RPC Application Server"
	a.AddFn("ping", func(in toolkit.M) *toolkit.Result {
		result := toolkit.NewResult()
		result.Data = "EACIIT RPC Application Server"
		return result
	}, false, "")

	a.AddFn("addsession", func(in toolkit.M) *toolkit.Result {
		referenceID := in.GetString("auth_referenceid")
		result := toolkit.NewResult()
		if referenceID == "" {
			return result.SetErrorTxt("Empty user provided")
		}

		//session, exist := a.sessions[referenceID]
		var session *Session

		for _, session = range a.sessions {
			if session.ReferenceID == referenceID && session.IsValid() && !a.AllowMultiLogin {
				return result.SetErrorTxt(referenceID + " already has active session on other connection")
			} else if session.ReferenceID == referenceID && !session.IsValid() && !a.AllowMultiLogin {
				delete(a.sessions, session.SessionID)
			}
		}
		session = a.RegisterSession(referenceID)
		//a.sessions[session.SessionID] = session
		//result.SetBytes(session, MarshallingMethod())
		result.Data = toolkit.M{}.Set("referenceid", session.SessionID).Set("secret", session.Secret).ToBytes("gob")
		return result
	}, true, "")

	a.AddFn("removesession", func(in toolkit.M) *toolkit.Result {
		result := toolkit.NewResult()
		referenceID := in.GetString("auth_referenceid")
		delete(a.sessions, referenceID)
		return result
	}, true, "session")

	//a.users = map[string]*User
	a.sessions = map[string]*Session{}
	a.listener = l
	go func() {
		rpc.Accept(l)
	}()
	return nil
}

func (a *Server) RegisterSession(referenceID string) *Session {
	s := NewSession(referenceID)
	s.Secret = toolkit.RandomString(32)
	s.SessionID = toolkit.RandomString(32)
	s.ReferenceID = referenceID
	a.sessions[s.SessionID] = s
	a.Log.Info(fmt.Sprintf("Registering new session [%s] for %s", s.SessionID, s.ReferenceID))
	return s
}

func (a *Server) AddUser(userid, password string) {
	user := new(User)
	user.ReferenceID = userid
	user.Secret = password
	if a.users == nil {
		a.users = map[string]*User{}
	}
	a.users[user.ReferenceID] = user
}

func (a *Server) AddFn(methodname string, fn RpcFn, needAuth bool, authType string) {
	var r *Rpc
	if a.rpcObject == nil {
		r = new(Rpc)
	} else {
		r = a.rpcObject
	}

	AddFntoRpc(r, a, methodname, fn, needAuth, authType)
	a.rpcObject = r
}

func (a *Server) Fn(fnName string) *RpcFnInfo {
	if a.rpcObject == nil {
		return nil
	}
	fnName = strings.ToLower(fnName)
	fn, exist := a.rpcObject.Fns[fnName]
	if !exist {
		return nil
	}
	return fn
}

func (a *Server) RegisterRPCFunctions(o interface{}) error {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		return errors.New("Invalid object for RPC Register")
	}
	objName := toolkit.TypeName(o)
	methodCount := t.NumMethod()
	for i := 0; i < methodCount; i++ {
		method := t.Method(i)
		mtype := method.Type
		methodName := strings.ToLower(method.Name)
		//fmt.Println("Evaluating " + toolkit.TypeName(o) + "." + methodName)

		//-- now check method signature
		if mtype.NumIn() == 2 && mtype.In(1).String() == "toolkit.M" {
			if mtype.NumOut() == 1 && mtype.Out(0).String() == "*toolkit.Result" {
				fmt.Println("Registering RPC Function " + objName + "." + methodName)
				a.AddFn(methodName, v.Method(i).Interface().(func(toolkit.M) *toolkit.Result), true, "session")
			}
		}
	}
	return nil
}

/*
func (a *Server) Serve() error {
	rpc.Accept(a.listener)
	return nil
}
*/

func (a *Server) Stop() error {
	a.listener.Close()
	a.Log.Info("Stopping service")
	return nil
}
