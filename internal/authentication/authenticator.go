package auth

// Our authentication service must implement Authenticator
// (we use the same Cassandra DB here for convenience;
// optimally, use different implementation for our authenticator service)
type Authenticator interface {
	Register() ([]byte, error)
	VerifyKey(string) (bool, error)
}
