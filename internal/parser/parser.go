package parser

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
	"go.dev.pztrn.name/hosts-translator/internal/models"
)

var ErrParserError = errors.New("hosts parser")

// Parser is a controlling structure for hosts file (both parsed and unparsed data).
type Parser struct {
	parsedHosts []models.Host
}

// NewParser creates new hosts file parsing controlling structure.
func NewParser() *Parser {
	p := &Parser{}
	p.initialize()

	return p
}

// Returns file data as bytes.
func (p *Parser) getFileData() ([]string, error) {
	// Before everything we should normalize file path.
	filePath := configuration.HostsFilePath

	// Replace possible "~" in the beginning as file reading function unable
	// to expand it. Also we should check only beginning because tilde is actually
	// a very valid directory or file name.
	if strings.HasPrefix(filePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		filePath = strings.Replace(filePath, "~", homeDir, 1)
	}

	// Get absolute file path.
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	log.Println("Reading file data from", absolutePath)

	// Read file data.
	rawData, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	data := strings.Split(string(rawData), "\n")

	return data, nil
}

// GetParsedData returns parsed data slice.
func (p *Parser) GetParsedData() []models.Host {
	return p.parsedHosts
}

// Initializes internal state of parser as well as CLI flags.
func (p *Parser) initialize() {
	p.parsedHosts = make([]models.Host, 0)
}

// Parse parses hosts file into internal representation.
func (p *Parser) Parse() error {
	log.Println("Starting hosts file parsing. File located at", configuration.HostsFilePath)

	data, err := p.getFileData()
	if err != nil {
		return err
	}

	for _, line := range data {
		// We should skip commented lines.
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Every line is a two-or-more-not-empty-strings. First string is always
		// an IP address.
		// Also there are a non-zero possibility that line will contain tabs, so as
		// very first action we should replace them with spaces.
		if strings.Contains(line, "\t") {
			line = strings.Replace(line, "\t", " ", -1)
		}

		lineSplitted := strings.Split(line, " ")

		// As one IP address can be bound to multiple domains we should take care
		// of that situation by creating multiple Host structures.
		var address string

		for _, lineData := range lineSplitted {
			// Also there might be a case when address placed first in line but
			// line itself has spaces in the beginning.
			if address == "" && lineData != "" {
				address = lineData

				continue
			}

			if lineData == "" {
				continue
			}

			domainToAdd := lineData

			if configuration.DomainPostfix != "" && !strings.HasSuffix(domainToAdd, configuration.DomainPostfix) {
				domainToAdd += "." + configuration.DomainPostfix
			}

			p.parsedHosts = append(p.parsedHosts, models.Host{
				Domain:  domainToAdd,
				Address: address,
			})
		}
	}

	log.Println("Got", len(p.parsedHosts), "domains from hosts file")

	// ToDo: hide under CLI parameter like '-debug'?
	// log.Printf("%+v\n", p.parsedHosts)

	return nil
}
