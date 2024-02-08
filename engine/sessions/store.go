package sessions

type Store interface {
	Set(*Session) error
	Get(sid string) (*Session, error)
	Purge(sid string) error
}
