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
	// Gets session or generates a new one determined by the following conditions.
	//
	// If r.Cookie[opts.Name] returns a cookie, it's value [opts.Value] is authenticated
	// and checked if it was legitemately signed by the backend [engine.SecureCookie] crypter.
	// When that's the case it will load the cookie values from the database.
	//
	// All cookies in simpa are encrypted and can be decrypted with [engine.SecureCookie]
	// Standard algorithms are AES_CTR and AES_GCM. Our API allows you to also define your own
	// crypter for customization.
	//
	// When no cookie can be found or seems to be not legitemately signed by backend
	// a new one is created instead. It wont be stored
	// unles the clients calls the Save() method
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
	New(r *http.Request, name string) *Session

	Save(session *Session) error

	Purge() error

	Get() *Session
}
