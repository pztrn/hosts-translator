package configuration

import (
	"flag"
	"fmt"
	"strings"
)

var (
	// PowerDNSAPIKey defines API key for PowerDNS HTTP API.
	PowerDNSAPIKey string
	// PowerDNSURI defines URL for PowerDNS HTTP API.
	PowerDNSURI string
)

func initializePowerDNS() {
	flag.StringVar(&PowerDNSAPIKey, "powerdns-api-key", "", "API key for PowerDNS HTTP API.")
	flag.StringVar(&PowerDNSURI, "powerdns-uri", "", "URI for PowerDNS API. Should be in 'proto://ADDR:PORT' form.")
}

func validatePowerDNS() error {
	if strings.ToLower(StorageToUse) != "powerdns" {
		return nil
	}

	// PowerDNS storage requires DomainSuffix to determine zone name to update.
	if DomainPostfix == "" {
		return fmt.Errorf("%w: domain postfix isn't filled which is required by PowerDNS storage", ErrConfigurationError)
	}

	if PowerDNSAPIKey == "" {
		return fmt.Errorf("%w: no PowerDNS API key was provided", ErrConfigurationError)
	}

	if PowerDNSURI == "" {
		return fmt.Errorf("%w: no PowerDNS HTTP API server URI provided", ErrConfigurationError)
	}

	// Hack: trim slashes in end.
	PowerDNSURI = strings.TrimRight(PowerDNSURI, "/")

	return nil
}
