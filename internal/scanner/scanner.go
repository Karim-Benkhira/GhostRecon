package scanner

import (
	"context"
	"net"
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


func (s *Scanner) ScanPort(ctx context.Context, port int) (*Result, error) {
	address := net.JoinHostPort(s.options.Target, string(port))
	conn, err := net.DialTimeout("tcp", address, s.options.Timeout)
	if err != nil {
		return &Result{
			Port:  port,
			State: "closed",
		}, nil
	}
	defer conn.Close()

	return &Result{
		Port:  port,
		State: "open",
	}, nil
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
