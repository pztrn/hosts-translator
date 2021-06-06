package powerdns

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
	"go.dev.pztrn.name/hosts-translator/internal/parser"
)

var ErrPowerDNSStorageError = errors.New("powerdns storage")

// PowerDNS is a controlling structure for PowerDNS storage.
type PowerDNS struct {
	parser *parser.Parser

	httpClient *http.Client
}

// NewPowerDNS creates new PowerDNS storage.
func NewPowerDNS(p *parser.Parser) *PowerDNS {
	s := &PowerDNS{
		parser: p,
	}

	s.initialize()

	return s
}

// Initializes storage internal state.
func (s *PowerDNS) initialize() {
	s.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
}

// Process processes parsed data.
func (s *PowerDNS) Process() error {
	log.Println("Processing parsed data into PowerDNS storage...")

	zoneData, err := s.getZoneData(configuration.DomainPostfix)
	if err != nil {
		return err
	}

	hostsParsedData := s.parser.GetParsedData()

	// Check what we should to create, update or delete.
	recordsToCreate := make([]RRSet, 0)
	recordsToDelete := make([]RRSet, 0)
	recordsToUpdate := make([]RRSet, 0)

	// First iteration - figure out what to create or update.
	for _, hostData := range hostsParsedData {
		var found bool

		for _, rrset := range zoneData.RRSets {
			// We're only for A or AAAA things here.
			if strings.ToUpper(rrset.Type) != "A" && strings.ToUpper(rrset.Type) != "AAAA" {
				continue
			}

			// ToDo: multiple addresses support somehow?
			if strings.TrimSuffix(rrset.Name, ".") == hostData.Domain {
				found = true

				if rrset.Records[0].Content != hostData.Address {
					recordsToUpdate = append(recordsToUpdate, RRSet{
						ChangeType: "REPLACE",
						Name:       rrset.Name,
						Records: []Record{
							{
								Content: hostData.Address,
							},
						},
						TTL:  300,
						Type: rrset.Type,
					})
				}

				break
			}
		}

		if !found {
			recordsToCreate = append(recordsToCreate, RRSet{
				ChangeType: "REPLACE",
				Name:       hostData.Domain + ".",
				Records: []Record{
					{
						Content: hostData.Address,
					},
				},
				TTL: 300,
				// ToDo: support for AAAA?
				Type: "A",
			})
		}
	}

	// Second iteration - figure out what to delete.
	for _, rrset := range zoneData.RRSets {
		// We're only for A or AAAA things here.
		if strings.ToUpper(rrset.Type) != "A" || strings.ToUpper(rrset.Type) != "AAAA" {
			continue
		}

		var found bool

		for _, hostData := range hostsParsedData {
			if strings.TrimSuffix(rrset.Name, ".") == hostData.Domain+"." {
				found = true

				break
			}
		}

		if !found {
			rrset.ChangeType = "DELETE"
			recordsToDelete = append(recordsToDelete, rrset)
		}
	}

	log.Println("Got", len(zoneData.RRSets), "RRSets in NS")
	log.Println("Got", len(recordsToCreate), "RRSets to create")
	log.Println("Got", len(recordsToUpdate), "RRSets to update")
	log.Println("Got", len(recordsToDelete), "RRSets to delete")

	recordsUnchanged := len(zoneData.RRSets) - len(recordsToDelete) - len(recordsToUpdate)
	log.Println("Got", recordsUnchanged, "RRSets unchanged")

	// ToDo: '-debug'?
	log.Printf("Got RRSets to create: %+v\n", recordsToCreate)
	// log.Printf("Got RRSets to update: %+v\n", recordsToUpdate)
	// log.Printf("Got RRSets to delete: %+v\n", recordsToDelete)

	s.updateZoneData(zoneData.Name, append(recordsToCreate, append(recordsToUpdate, recordsToDelete...)...))

	return nil
}
