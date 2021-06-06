package powerdns

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
)

// Updates zone data from PowerDNS.
func (s *PowerDNS) updateZoneData(zoneName string, RRSets []RRSet) error {
	log.Println("Updating zone data for domain", zoneName)

	zd := &zoneData{
		Name:   zoneName,
		RRSets: RRSets,
		Type:   "Zone",
	}

	url := strings.Join([]string{configuration.PowerDNSURI, "api", "v1", "servers", "localhost", "zones", strings.TrimSuffix(zd.Name, ".")}, "/")

	log.Println("URL:", url)

	zoneBytes, err := json.Marshal(zd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(zoneBytes))
	if err != nil {
		return err
	}

	bytesData, err := s.request(req)
	if err != nil {
		return err
	}

	log.Println("Got response:", string(bytesData))

	return nil
}
