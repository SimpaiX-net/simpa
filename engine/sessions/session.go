package sessions

import (
	"crypto/rand"
	"net/http"

	"github.com/SimpaiX-net/simpa/engine/crypt"
)

type Session struct {
	ID     string
	Values map[string]any
	Opts   *http.Cookie
	store  StoreI
	crypt  crypt.CrypterI
}

type SessionI interface {
	// Creates or loads session from the database
	New(r *http.Request, config *http.Cookie) (*Session, error)

	

	// generateres unique SID
	genSID() (string, error)
}

func (s *Session) New(r *http.Request, config *http.Cookie) (*Session, error) {
	// load cookie
	if c, err := r.Cookie(config.Name); err == nil {
		uid, err := s.crypt.Decrypt(c.Value)
		 if err != nil {
			return nil, err
		}

		return s.store.Get(uid)
	}

	sid, err := s.genSID()
	if err != nil {
		return nil, err
	}

	s.ID = sid
	s.Opts = config

	return s, nil
}

func (s *Session) genSID() (string, error) {
	var buff []byte
	for {
		buff = make([]byte, 32)
		if _, err := rand.Read(buff); err != nil {
			return "", err
		}

		if _, err := s.store.Get(string(buff)); err != nil {
			continue
		}

		break
	}
	return string(buff), nil
}
