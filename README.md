# GhostRecon - Advanced Reconnaissance Tool

## Overview
GoRecon is a powerful reconnaissance tool written in Go, designed for ethical hackers and security researchers. It provides fast and efficient information gathering capabilities about target systems and networks.

## Features
- 🔍 **Subdomain Enumeration**
  - DNS enumeration
  - Certificate transparency logs
  - Integration with various APIs (Shodan, VirusTotal)

- 🌐 **Port Scanning**
  - Fast concurrent port scanning using Goroutines
  - Support for both quick and deep scan modes
  - Common and custom port ranges

- 🔐 **Technology Detection**
  - Web technology stack fingerprinting
  - SSL/TLS certificate analysis
  - Header analysis
  - Service version detection

- 📊 **Multiple Operation Modes**
  - Passive mode (no direct target interaction)
  - Active mode (full reconnaissance)
  - API mode (for integration with other tools)

## Installation
```bash
go install github.com/yourusername/gorecon@latest
```

## Quick Start
```bash
# Basic scan
gorecon scan example.com

# Full scan with all features
gorecon scan --full example.com

# Passive mode scan
gorecon scan --passive example.com

# Export results to JSON
gorecon scan example.com --output results.json
```

## Configuration
- Default configuration file location: `~/.gorecon/config.yaml`
- API keys can be set via environment variables or config file
- Custom templates supported for output formatting

## Usage Examples
```bash
# Subdomain enumeration only
gorecon enum --subdomains example.com

# Port scanning with custom range
gorecon scan --ports 80,443,8000-8080 example.com

# Technology detection
gorecon fingerprint example.com
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## Security
- Please use this tool responsibly
- Always obtain proper authorization before scanning any systems
- Report security issues via GitHub security advisories

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- Thanks to the Go community
- Inspired by various open-source security tools