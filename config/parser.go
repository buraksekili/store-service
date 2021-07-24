package config

// ConfigParses represents an interface for reading environment
// variables for Store services.
type ConfigParser interface {
	Parse() error
	Address() (string, error)
}
