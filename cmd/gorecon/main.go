package main

import (
	"fmt"
	"os"

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
		fmt.Printf("Starting subdomain enumeration for: %s\n", target)

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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
