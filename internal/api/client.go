package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type APIClient struct {
	APIKey     string
	HTTPClient *http.Client
}


func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}


type ShodanClient struct {
	*APIClient
}


func NewShodanClient(apiKey string) *ShodanClient {
	return &ShodanClient{
		APIClient: NewAPIClient(apiKey),
	}
}


func (c *ShodanClient) SearchDomain(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.shodan.io/dns/domain/%s?key=%s", domain, c.APIKey)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Subdomains []string `json:"subdomains"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Subdomains, nil
}


type VirusTotalClient struct {
	*APIClient
}


func NewVirusTotalClient(apiKey string) *VirusTotalClient {
	return &VirusTotalClient{
		APIClient: NewAPIClient(apiKey),
	}
}


func (c *VirusTotalClient) SearchDomain(domain string) ([]string, error) {
	url := fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", c.APIKey, domain)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Subdomains []string `json:"subdomains"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Subdomains, nil
}
