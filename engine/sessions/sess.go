// TODO
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
	// Opts should contain cookie options, minimum required field is
	// [opts.Name]
	//
	// Loads session from datastore and syncs with client session cookie
	// when one is present or creates new one.
	//
	// For authentication and integrity auth it uses ``app.SecureCookie`` crypter.
	// Default crypters are AES CTR HMAC or AES GCM HMAC
	New(r *http.Request, opts http.Cookie) (*Session, error)

	// Saves the session to the database.
	// It also overwrites client cookie
	//
	// The client cookie holds cookie.Name = sessionName
	// cookie.Value = [encrypted value] -> using ``app.SecureCookie``
	Save() error

	// Gets session value [key]
	GetVal(key string) (interface{}, error)

	// Sets session value [key]=val
	SetVal(key string, val interface{}) error

	// Deletes session value [key]=val
	DeleteVal(key string) error

	// Deletes all session values
	DeleteValAll() error

	// Destroys the session both from database and client
	Destroy() error
}

// Store implementation for sessions.
// Uses app.SecureCookie crypter to crypt. Standard crypters use AES_CTR HMACor AES_GCM HMAC.
type SessionStore interface {
	// Loads session from the session cookie sent form the client or
	// generates a new one.
	//
	// When loading, the client cookie is authenticated and then loaded from data store.
	// For setting any cookies Simpa does use ``engine.SecureCookie``.
	// ``engine.SecureCookie`` contains a crypter model.
	//
	// Defaults to AES GCM HMAC or AES CTR HMAC crypter near choice.
	//
	// To summarize when loading cookies its authenticity and integrity
	// is authenticated first before trying to load.
	//
	// Otherwise it creates a new session
	New(name string) *Session

	Save(session *Session) error

	Purge(id string) error

	Get(key string) interface{}
}
