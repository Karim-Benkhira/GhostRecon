package utils

import (
	"encoding/json"
	"os"
)


type ExportResult struct {
	Target     string              `json:"target"`
	Timestamp  string              `json:"timestamp"`
	Subdomains []SubdomainExport   `json:"subdomains,omitempty"`
	DNSRecords map[string][]string `json:"dns_records,omitempty"`
}


type SubdomainExport struct {
	Name   string   `json:"name"`
	IPs    []string `json:"ips"`
	Status string   `json:"status"`
}


func ExportToJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}
