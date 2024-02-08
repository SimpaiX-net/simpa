package sessions

import (
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"simpaix.net/simpa/engine/crypt"
)

var ErrCipher = errors.New("simpa-sessions: SID is not signed by backend, but for compability created a new session")

// Session Cookie cofig
type Config struct {
	Name     string        `bson:"name" json:"name"`
	Value    string        `bson:"value" json:"value"`
	MaxAge   int           `bson:"maxage" json:"maxage"`
	Expires  time.Time     `bson:"expires" json:"expires"`
	Secure   bool          `bson:"secure" json:"secure"`
	HttpOnly bool          `bson:"httponly" json:"httponly"`
	SameSite http.SameSite `bson:"samesite" json:"samesite"`
}

// Session object
type Session struct {
	ID     string                 `bson:"_id" json:"_id"`
	Values map[string]interface{} `bson:"values" json:"values"` // when marshaled should be encrypted through engine's crypter, when unmarshalled should decrypt and parse json
	Opts   *Config                `bson:"options" json:"options"`
	store  Store
	crypt  crypt.CrypterI
}

type SessionI interface {
	// Creates or loads session from the database
	New(r *http.Request, config http.Cookie) (*Session, error)

	// generateres unique SID
	genSID() (string, error)

	// saves session back to store
	Save() error

	// returns SID
	SID() string

	// sets key to map
	Set(key string, val interface{})

	// gets ley from map
	Get(key string) interface{}

	// sets store driver
	SetStore(Store)

	// sets crypter, should set it from engine
	SetCrypter(crypt.CrypterI)
}

func (s *Session) New(r *http.Request, config *Config) (*Session, error) {
	var err error
	if c, noCookie := r.Cookie(config.Name); noCookie == nil {
		// cookie on user's browser found
		sid, err := s.crypt.Decrypt(c.Value)
		if err == nil {
			sess, err := s.store.Get(sid)
			if err != nil {
				return nil, err
			}

			s.ID = sess.ID
			s.Values = sess.Values
			s.Opts = sess.Opts

			return s, err
		}

	}

	gid, err := s.genSID()
	if err != nil {
		return nil, err
	}

	s.ID = gid
	s.Opts = config

	return s, err
}

func (s *Session) Save(w http.ResponseWriter) error {
	if err := s.store.Set(s); err != nil {
		return err
	}

	sid, err := s.crypt.Encrypt(s.ID)
	if err != nil {
		return err
	}

	s.Opts.Value = sid

	http.SetCookie(w, &http.Cookie{
		Name:     s.Opts.Name,
		Value:    s.Opts.Value,
		MaxAge:   s.Opts.MaxAge,
		Expires:  s.Opts.Expires,
		Secure:   s.Opts.Secure,
		HttpOnly: s.Opts.HttpOnly,
		SameSite: s.Opts.SameSite,
	})
	return nil
}

func (s *Session) SID() string {
	return s.ID
}

func (s *Session) genSID() (string, error) {
	var sid string
	for {
		buff := make([]byte, 12)
		if _, err := rand.Read(buff); err != nil {
			return sid, err
		}

		sid = primitive.ObjectID(buff).Hex()
		if _, err := s.store.Get(sid); err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
			return sid, err
		}
	}

	return sid, nil
}

func (s *Session) SetStore(store Store) {
	s.store = store
}

func (s *Session) SetCrypter(crypt crypt.CrypterI) {
	s.crypt = crypt
}

func (s *Session) Set(key string, val interface{}) {
	s.Values[key] = val
}

func (s *Session) Get(key string) interface{} {
	return s.Values[key]
}
