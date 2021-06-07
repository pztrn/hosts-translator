package powerdns

type zoneData struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	RRSets []RRSet `json:"rrsets"`
	Serial int64   `json:"serial"`
}

type RRSet struct {
	ChangeType string   `json:"changetype"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Records    []Record `json:"records"`
	TTL        int64    `json:"ttl"`
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
