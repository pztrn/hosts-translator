package powerdns

import (
	"io"
	"net/http"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
)

// Executes request to PowerDNS server and returns data or error.
func (s *PowerDNS) request(req *http.Request) ([]byte, error) {
	req.Header.Add("X-API-Key", configuration.PowerDNSAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
