package sessions

type StoreI interface {
	Set(*Session) error
	Get(sid string) (*Session, error)
	Purge(sid string) error
}
