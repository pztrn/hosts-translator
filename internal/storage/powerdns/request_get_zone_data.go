package powerdns

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
)

// Gets zone data from PowerDNS.
func (s *PowerDNS) getZoneData(zoneName string) (*zoneData, error) {
	log.Println("Getting zone data for domain", zoneName)

	url := strings.Join([]string{configuration.PowerDNSURI, "api", "v1", "servers", "localhost", "zones", zoneName}, "/")

	log.Println("URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bytesData, err := s.request(req)
	if err != nil {
		return nil, err
	}

	zd := &zoneData{}

	if err := json.Unmarshal(bytesData, zd); err != nil {
		return nil, err
	}

	return zd, nil
}
