package recon

import (
	"context"
	"fmt"
	"net"

	"github.com/Karim-Benkhira/GhostRecon/internal/api"
	"github.com/Karim-Benkhira/GhostRecon/internal/utils"
)

type EnumOptions struct {
	Domain     string
	Recursive  bool
	MaxDepth   int
	Concurrent int
	ShodanKey  string
	VTotalKey  string
	UseAPIs    bool
}

type SubdomainResult struct {
	Subdomain string
	IPs       []string
	Status    string
}

type Enumerator struct {
	options EnumOptions
}

func NewEnumerator(opts EnumOptions) *Enumerator {
	return &Enumerator{
		options: opts,
	}
}


var CommonWordlist = []string{
	"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop",
	"ns1", "ns2", "dns", "dns1", "dns2", "ns", "dev", "staging",
	"api", "admin", "mx", "ssh", "vpn", "web", "test", "portal",
}


var DNSRecordTypes = []string{"A", "AAAA", "CNAME", "MX", "NS", "TXT"}


func (e *Enumerator) EnumerateBasic(ctx context.Context) []SubdomainResult {
	var results []SubdomainResult

	for _, prefix := range CommonWordlist {
		select {
		case <-ctx.Done():
			return results
		default:
			subdomain := prefix + "." + e.options.Domain
			ips, err := e.ResolveDomain(subdomain)

			if err == nil {
				results = append(results, SubdomainResult{
					Subdomain: subdomain,
					IPs:       ips,
					Status:    "active",
				})
			}
		}
	}

	return results
}


func (e *Enumerator) Start(ctx context.Context) chan SubdomainResult {
	results := make(chan SubdomainResult)

	go func() {
		defer close(results)

		
		semaphore := make(chan struct{}, e.options.Concurrent)

		
		for _, prefix := range CommonWordlist {
			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
				go func(p string) {
					defer func() { <-semaphore }()

					subdomain := p + "." + e.options.Domain
					ips, err := e.ResolveDomain(subdomain)

					if err == nil {
						results <- SubdomainResult{
							Subdomain: subdomain,
							IPs:       ips,
							Status:    "active",
						}
					}
				}(prefix)
			}
		}
	}()

	return results
}


func (e *Enumerator) CheckDNSRecords(domain string) map[string][]string {
	records := make(map[string][]string)

	for _, recordType := range DNSRecordTypes {
		
		if ips, err := e.ResolveDomain(domain); err == nil {
			records[recordType] = ips
		}
	}

	return records
}


func (e *Enumerator) ResolveDomain(domain string) ([]string, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		var errType utils.ErrorType
		if dnsErr, ok := err.(*net.DNSError); ok {
			errType = utils.DNSError
			if dnsErr.Timeout() {
				errType = utils.TimeoutError
			}
		} else {
			errType = utils.NetworkError
		}
		return nil, utils.NewError(errType,
			fmt.Sprintf("failed to resolve domain %s", domain),
			err)
	}
	return ips, nil
}


func (e *Enumerator) EnumerateWithAPIs(ctx context.Context) ([]SubdomainResult, error) {
	var results []SubdomainResult
	var apiErrors []error

	if e.options.UseAPIs {
		if e.options.ShodanKey != "" {
			shodan := api.NewShodanClient(e.options.ShodanKey)
			subdomains, err := shodan.SearchDomain(e.options.Domain)
			if err != nil {
				apiErrors = append(apiErrors, utils.NewError(utils.APIError,
					"Shodan API error", err))
			} else {
				for _, sub := range subdomains {
					if ips, err := e.ResolveDomain(sub); err == nil {
						results = append(results, SubdomainResult{
							Subdomain: sub,
							IPs:       ips,
							Status:    "active",
						})
					}
				}
			}
		}

		if e.options.VTotalKey != "" {
			vt := api.NewVirusTotalClient(e.options.VTotalKey)
			subdomains, err := vt.SearchDomain(e.options.Domain)
			if err != nil {
				apiErrors = append(apiErrors, utils.NewError(utils.APIError,
					"VirusTotal API error", err))
			} else {
				for _, sub := range subdomains {
					if ips, err := e.ResolveDomain(sub); err == nil {
						results = append(results, SubdomainResult{
							Subdomain: sub,
							IPs:       ips,
							Status:    "active",
						})
					}
				}
			}
		}
	}

	if len(apiErrors) > 0 {
		return results, fmt.Errorf("API errors occurred: %v", apiErrors)
	}
	return results, nil
}
