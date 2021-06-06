package models

// Host represents structure of single host that was parsed from hosts file.
type Host struct {
	Domain  string
	Address string
}
