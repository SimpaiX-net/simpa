package sessions

import "net/http"

// Represents a web session
type Session struct {
	ID      string         `json:"id"`      // Session ID
	Values  map[string]any `json:"values"`  // Decrypted session context object
	Options *http.Cookie   `json:"options"` // Cookie options
	store   SessionStore   // Store wrapper
}

type SessionI interface {
	New(r *http.Request, opts http.Cookie) (*Session, error)

	// Saves the session to the database.
	// It also rewrites and sends it to the client
	Save() error

	// Gets session value [key]
	GetVal(key string) (interface{}, error)

	// Sets session value [key]=val
	SetVal(key string, val interface{}) error

	// Deletes session value [key]=val
	DeleteVal(key string) error

	// Destroys the session both from database and client
	Destroy() error
}

// Store implementation for sessions.
// Uses app.SecureCookie crypter to crypt. Standard crypters use AES_CTR or AES_GCM.
type SessionStore interface {
	New(id string) *Session

	Save(session *Session) error

	Purge() error

	Get() *Session
}
