package display

import (
	"fmt"
	"strings"
	"time"
)


const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)


type ResultPrinter struct {
	ShowTimestamp bool
	Verbose       bool
}


func NewPrinter(showTimestamp, verbose bool) *ResultPrinter {
	return &ResultPrinter{
		ShowTimestamp: showTimestamp,
		Verbose:       verbose,
	}
}


func (p *ResultPrinter) PrintHeader(text string) {
	fmt.Printf("\n%s=== %s ===%s\n", ColorCyan, text, ColorReset)
}


func (p *ResultPrinter) PrintSubdomain(subdomain string, ips []string, status string) {
	timestamp := ""
	if p.ShowTimestamp {
		timestamp = time.Now().Format("15:04:05 ")
	}

	statusColor := ColorGreen
	if status != "active" {
		statusColor = ColorRed
	}

	fmt.Printf("%s%s[%s%s%s] %s%s%s\n",
		ColorBlue, timestamp,
		statusColor, status, ColorBlue,
		subdomain, ColorReset,
		p.formatIPs(ips),
	)
}


func (p *ResultPrinter) PrintDNSRecord(recordType string, values []string) {
	if len(values) > 0 {
		fmt.Printf("%s%s Records:%s %s\n",
			ColorPurple, recordType, ColorReset,
			strings.Join(values, ", "),
		)
	}
}


func (p *ResultPrinter) PrintSummary(totalSubdomains int, activeSubdomains int, duration time.Duration) {
	p.PrintHeader("Scan Summary")
	fmt.Printf("\n%sTotal Subdomains:%s %d\n", ColorYellow, ColorReset, totalSubdomains)
	fmt.Printf("%sActive Subdomains:%s %d\n", ColorYellow, ColorReset, activeSubdomains)
	fmt.Printf("%sScan Duration:%s %s\n\n", ColorYellow, ColorReset, duration)
}

func (p *ResultPrinter) formatIPs(ips []string) string {
	if len(ips) == 0 {
		return ""
	}
	return fmt.Sprintf(" (%s%s%s)", ColorYellow, strings.Join(ips, ", "), ColorReset)
}
