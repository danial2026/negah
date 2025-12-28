package scanner

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ScanFeature helps us keep track of what each nmap command does
type ScanFeature struct {
	ID          int
	Name        string
	Description string
	Command     string // The actual flag or nmap script to run
	Sudo        bool   // Some stuff needs root/admin privileges
}

// CheckNmap just checks if the user actually has nmap installed
// If it's not in the PATH, we should probably warn them
func CheckNmap() bool {
	_, err := exec.LookPath("nmap")
	return err == nil
}

// ExecuteCommand is a helper to run commands and show the progress directly
func ExecuteCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Just logging what we're about to run for transparency
	fmt.Printf("\n[Running] %s %s\n\n", name, strings.Join(args, " "))
	return cmd.Run()
}

// ExecuteCommandWithOutput runs a command and captures its output
func ExecuteCommandWithOutput(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Build command info
	cmdInfo := fmt.Sprintf("[Running] %s %s\n\n", name, strings.Join(args, " "))
	
	err := cmd.Run()
	
	// Combine command info with output
	output := cmdInfo + stdout.String()
	if stderr.Len() > 0 {
		output += "\n[Errors]\n" + stderr.String()
	}
	
	return output, err
}

// GetPublicIP reaches out to ipapi.co to get our ip and location
func GetPublicIP() {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		fmt.Printf("Couldn't grab public IP info: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n--- Your Public IP Info ---\n%s\n", string(body))
}

// GetPublicIPWithOutput returns public IP info as a string
func GetPublicIPWithOutput() (string, error) {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return "", fmt.Errorf("couldn't grab public IP info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("--- Your Public IP Info ---\n%s\n", string(body)), nil
}

// GetLocalInfo checks our own interfaces to see what's happening locally
func GetLocalInfo() {
	fmt.Println("\n--- Local Interface Details ---")
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Had trouble reading interfaces: %v\n", err)
		return
	}

	for _, i := range interfaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			fmt.Printf("Interface: %s | Address: %s\n", i.Name, addr.String())
		}
	}

	fmt.Printf("\nRunning on %s (%s)\n", runtime.GOOS, runtime.GOARCH)
}

// GetLocalInfoWithOutput returns local interface info as a string
func GetLocalInfoWithOutput() (string, error) {
	var output strings.Builder
	output.WriteString("--- Local Interface Details ---\n")
	
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("had trouble reading interfaces: %v", err)
	}

	for _, i := range interfaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			output.WriteString(fmt.Sprintf("Interface: %s | Address: %s\n", i.Name, addr.String()))
		}
	}

	output.WriteString(fmt.Sprintf("\nRunning on %s (%s)\n", runtime.GOOS, runtime.GOARCH))
	return output.String(), nil
}

// GetFeatures returns all 35 tools we've packed into this script
func GetFeatures() []ScanFeature {
	return []ScanFeature{
		{1, "Local Discovery", "Simple ping sweep of the network", "-sn", false},
		{2, "Hell Scan (Full)", "Checks every single port (1-65535)", "-p-", false},
		{3, "Quick Check", "Top 100 common ports only", "--top-ports 100", false},
		{4, "Custom Range", "You pick the ports", "-p", false},
		{5, "Service Versions", "What software is actually running?", "-sV", false},
		{6, "OS Detection", "Guesses what OS the target has", "-O", true},
		{7, "Total Takedown", "Aggressive scan with everything enabled", "-A", true},
		{8, "Vuln Finder", "Standard vulnerability scripts", "--script vuln", false},
		{9, "UDP Hunt", "Scanning for those tricky UDP ports", "-sU", true},
		{10, "Firewall Proofing", "ACK scan to see if there's a firewall", "-sA", false},
		{11, "Path Tracer", "See the hops to the target", "--traceroute", false},
		{12, "My Public IP", "Geo-info and public IP details", "INTERNAL_PUBLIC_IP", false},
		{13, "Local Network Info", "Local IPs, DNS, and interface list", "INTERNAL_LOCAL_INFO", false},
		{14, "Fast Mode", "Default nmap but faster", "-F", false},
		{15, "Ping Only", "Check if the host is even alive", "-sP", false},
		{16, "Web Titles", "Grabs HTTP headers and page titles", "--script http-title,http-headers", false},
		{17, "SSL/TLS Check", "Certificates and cipher suites audit", "--script ssl-cert,ssl-enum-ciphers", false},
		{18, "SMB OS Guess", "Detailed OS info via SMB", "--script smb-os-discovery", false},
		{19, "DNS Brute", "Try to guess subdomains", "--script dns-brute", false},
		{20, "SSH Audit", "Check for weak SSH ciphers", "--script ssh2-enum-algos", false},
		{21, "DB Hunt", "Search for MySQL, Postgres, Redis, etc.", "-p 3306,5432,6379,27017,1433", false},
		{22, "Banner Grabber", "Grab service banners for identification", "-sV --script banner", false},
		{23, "No DNS Resolving", "Fast scan without resolving names", "-n", false},
		{24, "Sneaky Scan", "Slow timing (T2) to avoid detection", "-T2", false},
		{25, "Protocol Scan", "See what IP protocols are supported", "-sO", true},
		{26, "SCTP Init Scan", "Specific for SCTP protocol", "-sY", false},
		{27, "FIN Scan", "Stealth scan (FIN packet)", "-sF", false},
		{28, "XMAS Scan", "Stealth scan (FIN, PSH, URG lights)", "-sX", false},
		{29, "Null Scan", "Stealth scan (no flags set)", "-sN", false},
		{30, "Fragmented Test", "Firewall test using fragmented packets", "-f", false},
		{31, "MTU Test", "Test firewall with specific MTU sizes", "--mtu 24", false},
		{32, "Source Port 53", "Spoof as DNS traffic to bypass filters", "--source-port 53", false},
		{33, "Bad Checksum", "Test if packets are dropped by firewall", "--badsum", false},
		{34, "Heartbleed Trip", "Check for the classic OpenSSL bug", "--script ssl-heartbleed", false},
		{35, "Whois Lookup", "Simple whois info for the domain", "INTERNAL_WHOIS", false},
	}
}

// RunScan is the main engine. It handles both internal and external (nmap) tools
func RunScan(feature ScanFeature, target string) {
	switch feature.Command {
	case "INTERNAL_PUBLIC_IP":
		GetPublicIP()
		return
	case "INTERNAL_LOCAL_INFO":
		GetLocalInfo()
		return
	case "INTERNAL_WHOIS":
		ExecuteCommand("whois", target)
		return
	}

	// Figure out if we need sudo
	args := []string{}
	if feature.Sudo && runtime.GOOS != "windows" {
		args = append(args, "sudo", "nmap")
	} else {
		args = append(args, "nmap")
	}

	// Splitting the command string into parts nmap understands
	cmdParts := strings.Fields(feature.Command)
	args = append(args, cmdParts...)

	// Add the IP or domain at the end
	if target != "" {
		args = append(args, target)
	}

	err := ExecuteCommand(args[0], args[1:]...)
	if err != nil {
		fmt.Printf("\nSomething went wrong during the scan: %v\n", err)
	}
}

// RunScanWithOutput is like RunScan but captures and returns the output
func RunScanWithOutput(feature ScanFeature, target string) (string, error) {
	switch feature.Command {
	case "INTERNAL_PUBLIC_IP":
		return GetPublicIPWithOutput()
	case "INTERNAL_LOCAL_INFO":
		return GetLocalInfoWithOutput()
	case "INTERNAL_WHOIS":
		return ExecuteCommandWithOutput("whois", target)
	}

	// Figure out if we need sudo
	args := []string{}
	cmdName := ""
	if feature.Sudo && runtime.GOOS != "windows" {
		cmdName = "sudo"
		args = append(args, "nmap")
	} else {
		cmdName = "nmap"
	}

	// Splitting the command string into parts nmap understands
	cmdParts := strings.Fields(feature.Command)
	args = append(args, cmdParts...)

	// Add the IP or domain at the end
	if target != "" {
		args = append(args, target)
	}

	return ExecuteCommandWithOutput(cmdName, args...)
}
