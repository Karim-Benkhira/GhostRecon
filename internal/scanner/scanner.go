package scanner

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)


type ScanOptions struct {
	Target      string
	Ports       []int
	Timeout     time.Duration
	Concurrency int
}


type Result struct {
	Port    int
	State   string
	Service string
}


type Scanner struct {
	options ScanOptions
}


func NewScanner(opts ScanOptions) *Scanner {
	return &Scanner{
		options: opts,
	}
}


var CommonPorts = []int{
	20, 21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 
	993, 995, 1723, 3306, 3389, 5900, 8080,
}

func ParsePortRange(portsStr string) ([]int, error) {
	if portsStr == "" {
		return CommonPorts, nil
	}

	var ports []int
	ranges := strings.Split(portsStr, ",")
	
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}
			
			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}
			
			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			
			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
	}
	
	return ports, nil
}

func (s *Scanner) ScanPort(ctx context.Context, port int) (*Result, error) {
	address := net.JoinHostPort(s.options.Target, strconv.Itoa(port))
	
	var d net.Dialer
	d.Timeout = s.options.Timeout
	
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return &Result{
			Port:    port,
			State:   "closed",
			Service: "unknown",
		}, nil
	}
	defer conn.Close()

	service := DetectService(port)
	
	return &Result{
		Port:    port,
		State:   "open",
		Service: service,
	}, nil
}

func DetectService(port int) string {
	commonServices := map[int]string{
		20: "FTP-data", 21: "FTP", 22: "SSH", 23: "Telnet",
		25: "SMTP", 53: "DNS", 80: "HTTP", 110: "POP3",
		143: "IMAP", 443: "HTTPS", 3306: "MySQL", 3389: "RDP",
		5432: "PostgreSQL", 8080: "HTTP-Proxy",
	}

	if service, exists := commonServices[port]; exists {
		return service
	}
	
	return "unknown"
}

func (s *Scanner) Start(ctx context.Context) chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		semaphore := make(chan struct{}, s.options.Concurrency)

		for _, port := range s.options.Ports {
			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
				go func(p int) {
					defer func() { <-semaphore }()

					result, _ := s.ScanPort(ctx, p)
					if result != nil {
						results <- *result
					}
				}(port)
			}
		}
	}()

	return results
}
