// Package config: the configuration API
package config

type Variable struct {
	Name  string
	Value string
}

type Configuration struct {
	Source    string
	Name      string
	Variables []Variable
	SSHKeys   []string
}
