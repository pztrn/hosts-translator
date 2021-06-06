package powerdns

type zoneData struct {
	Name   string  `json:"name"`
	RRSets []RRSet `json:"rrsets"`
	Serial int64   `json:"serial"`
	Type   string  `json:"type"`
}

type RRSet struct {
	ChangeType string   `json:"changetype"`
	Name       string   `json:"name"`
	Records    []Record `json:"records"`
	TTL        int64    `json:"ttl"`
	Type       string   `json:"type"`
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
