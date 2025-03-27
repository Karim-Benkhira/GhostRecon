package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Karim-Benkhira/GhostRecon/internal/display"
	"github.com/Karim-Benkhira/GhostRecon/internal/recon"
	"github.com/Karim-Benkhira/GhostRecon/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gorecon",
	Short: "GoRecon - Advanced Reconnaissance Tool",
	Long: `GoRecon is a powerful reconnaissance tool written in Go,
designed for ethical hackers and security researchers.
It provides fast and efficient information gathering capabilities.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Perform a reconnaissance scan on the target",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		fmt.Printf("Starting scan on target: %s\n", target)

	},
}

var enumCmd = &cobra.Command{
	Use:   "enum [target]",
	Short: "Enumerate subdomains for the target",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		startTime := time.Now()

		logger := utils.NewLogger(utils.INFO)
		printer := display.NewPrinter(true, true)

		
		var errors []error

		printer.PrintHeader("Starting Subdomain Enumeration")
		logger.Log(utils.INFO, "Target: %s", target)

		
		recursive, _ := cmd.Flags().GetBool("recursive")
		depth, _ := cmd.Flags().GetInt("depth")
		if depth < 1 {
			errors = append(errors, utils.NewError(utils.ConfigError,
				"invalid depth value", nil))
			depth = 1
		}

		
		opts := recon.EnumOptions{
			Domain:     target,
			Recursive:  recursive,
			MaxDepth:   depth,
			Concurrent: 10,
		}

		
		if !strings.Contains(target, ".") {
			errors = append(errors, utils.NewError(utils.ConfigError,
				"invalid domain format", nil))
			return
		}

		
		enumerator := recon.NewEnumerator(opts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		
		exportData := utils.ExportResult{
			Target:    target,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		
		results := enumerator.Start(ctx)
		var subdomains []utils.SubdomainExport

		
		var totalSubdomains, activeSubdomains int
		for result := range results {
			totalSubdomains++
			if result.Status == "active" {
				activeSubdomains++
			}

			subdomains = append(subdomains, utils.SubdomainExport{
				Name:   result.Subdomain,
				IPs:    result.IPs,
				Status: result.Status,
			})

			printer.PrintSubdomain(
				result.Subdomain,
				result.IPs,
				result.Status,
			)
		}

		
		printer.PrintHeader("DNS Records")
		records := enumerator.CheckDNSRecords(target)
		for recordType, values := range records {
			printer.PrintDNSRecord(recordType, values)
		}

		
		duration := time.Since(startTime)
		printer.PrintSummary(totalSubdomains, activeSubdomains, duration)

		
		exportData.Subdomains = subdomains
		exportData.DNSRecords = records

		
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile != "" {
			logger.Log(utils.INFO, "Exporting results to: %s", outputFile)
			if err := utils.ExportToJSON(outputFile, exportData); err != nil {
				logger.Log(utils.ERROR, "Failed to export results: %v", err)
			}
		}

		
		if len(errors) > 0 {
			printer.PrintHeader("Errors Occurred")
			for _, err := range errors {
				if reconErr, ok := err.(*utils.ReconError); ok {
					logger.Log(utils.ERROR, "[%v] %s", reconErr.Type, reconErr.Error())
				} else {
					logger.Log(utils.ERROR, "%v", err)
				}
			}
		}
	},
}

func init() {

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(enumCmd)

	scanCmd.Flags().Bool("passive", false, "Run in passive mode")
	scanCmd.Flags().Bool("full", false, "Perform a full scan")
	scanCmd.Flags().String("output", "", "Output file path")
	scanCmd.Flags().String("ports", "", "Custom port range (e.g., 80,443,8000-8080)")

	enumCmd.Flags().Bool("recursive", false, "Perform recursive enumeration")
	enumCmd.Flags().Int("depth", 1, "Recursion depth for enumeration")
	enumCmd.Flags().String("output", "", "Output file path for JSON results")
	enumCmd.Flags().String("shodan-key", "", "Shodan API key")
	enumCmd.Flags().String("vt-key", "", "VirusTotal API key")
	enumCmd.Flags().Bool("use-apis", false, "Use external APIs for enumeration")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
