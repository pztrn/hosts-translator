package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"go.dev.pztrn.name/hosts-translator/internal/configuration"
	"go.dev.pztrn.name/hosts-translator/internal/parser"
	"go.dev.pztrn.name/hosts-translator/internal/storage"
	"go.dev.pztrn.name/hosts-translator/internal/storage/powerdns"
)

func main() {
	log.Println("Starting hosts file translator...")

	configuration.Initialize()
	configuration.Parse()

	if err := configuration.Validate(); err != nil {
		log.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	p := parser.NewParser()

	var s storage.Interface

	switch strings.ToLower(configuration.StorageToUse) {
	case "powerdns":
		s = powerdns.NewPowerDNS(p)
	}

	if err := p.Parse(); err != nil {
		panic(err)
	}

	if err := s.Process(); err != nil {
		panic(err)
	}
}
