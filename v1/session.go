package appserver

import (
	"github.com/eaciit/toolkit"
	"time"
)

type Session struct {
	UserID   string
	Secret   string
	Created  time.Time
	ExpireOn time.Time
}

var _defaultSessionLifetime time.Duration

func SetSesionLifetime(t time.Duration) {
	_defaultSessionLifetime = t
}

func SessionLifetime() time.Duration {
	if _defaultSessionLifetime == 0 {
		_defaultSessionLifetime = 90 * time.Minute
	}
	return _defaultSessionLifetime
}

func NewSession(userid string) *Session {
	s := new(Session)
	s.UserID = userid
	s.Created = time.Now()
	s.ExpireOn = s.Created.Add(SessionLifetime())
	s.Secret = toolkit.GenerateRandomString("", 32)
	return s
}

func (s *Session) IsValid() bool {
	return time.Now().Before(s.ExpireOn)
}
