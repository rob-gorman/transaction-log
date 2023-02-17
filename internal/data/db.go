package data

// DBAccessor is the interface any database implementations must satisfy
type DBAccessor interface {
	InsertEvent([]byte) error
	SelectRowsByField(string, string) ([]byte, error)
}
