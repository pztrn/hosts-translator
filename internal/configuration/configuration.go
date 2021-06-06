package configuration

import (
	"errors"
	"flag"
	"fmt"
)

var (
	// ErrConfigurationError can be used to determine if received error was
	// about invalid configuration.
	ErrConfigurationError = errors.New("configuration")

	// DomainPostfix is a string that will be added to parsed domain if not already
	// present. E.g. with DomainPostfix = "example.com" for parsed domain "host-1"
	// it will be added and parsed domain will be in form "host-1.example.com", but
	// for parsed domain "host-2.example.com" it won't be added as it is already
	// present in the very end.
	DomainPostfix string
	// HostsFilePath defines a path from which hosts file will be used for parsing.
	HostsFilePath string
	// StorageToUse defines storage translator will use for updating data.
	StorageToUse string
)

// Initialize initializes configuration subsystem.
func Initialize() {
	flag.StringVar(&DomainPostfix, "domain-postfix", "", "Postfix to append to domain. Some storages requires this parameter to be filled.")
	flag.StringVar(&HostsFilePath, "hosts-file", "", "Path to hosts file to parse.")
	flag.StringVar(&StorageToUse, "storage", "", "Storage to use. Currently supported: 'powerdns'.")

	initializePowerDNS()
}

func Parse() {
	flag.Parse()
}

// Validate validates configuration data and returns error is something isn't right.
func Validate() error {
	if HostsFilePath == "" {
		return fmt.Errorf("%w: empty hosts file path", ErrConfigurationError)
	}

	if StorageToUse == "" {
		return fmt.Errorf("%w: no storage name was provided", ErrConfigurationError)
	}

	if err := validatePowerDNS(); err != nil {
		return err
	}

	return nil
}
