package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	runTool()
}

func runTool() {
	// Print the banner
	fmt.Println(`

███╗   ██╗██╗   ██╗███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗ 
████╗  ██║██║   ██║██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
██╔██╗ ██║██║   ██║███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██║╚██╗██║╚██╗ ██╔╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║ ╚████║ ╚████╔╝ ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═══╝  ╚═══╝  ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
                                                                               
                                                                                                                                                                                                                                              

	   Network Scanner By Dnyanesh & Ayush
	`)

	// Ask for target IP or domain name
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target IP or domain to scan: ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	// Display menu and get user selection
	fmt.Println("Select option from the menu (you can select multiple options separated by commas):")
	fmt.Println("1) Port Scan")
	fmt.Println("2) Host Discovery")
	fmt.Println("3) DNS Enumeration")
	fmt.Println("4) Service Identification")
	fmt.Println("5) Run Default Scripts")
	fmt.Println("6) Specify Port Range")
	fmt.Println("7) Service Version Identification")
	fmt.Println("8) OS Detection")
	fmt.Print("Enter your choice(s): ")
	choices, _ := reader.ReadString('\n')
	choices = strings.TrimSpace(choices)

	// Split choices into individual options
	optionList := strings.Split(choices, ",")

	// Execute selected options
	for _, option := range optionList {
		switch option {
		case "1":
			// Perform port scan
			fmt.Println("Starting port scan for", target)
			scanPorts(target)
		case "2":
			fmt.Print("Enter IP with subnet mask (e.g., 192.168.1.1/24): ")
			ipMask, _ := reader.ReadString('\n')
			ipMask = strings.TrimSpace(ipMask)
			runCommand("sh", "-c", fmt.Sprintf("nmap -sn %s | sed 's/Nmap/NVSCANNER/g' | awk '{gsub(\"Nmap\",\"NVSCANNER\")}1' | sed 's/https:\\/\\/nmap.org//g'", ipMask))
		case "3":
			runCommand("dnsrecon", "-d", target, "-z")
		case "4":
			runCommand("sh", "-c", fmt.Sprintf("nmap %s | sed 's/Nmap/NVSCANNER/g' | awk '{gsub(\"Nmap\",\"NVSCANNER\")}1' |  awk '/^PORT/{flag=1} flag; /^Nmap done/{flag=0}'", target))
		case "5":
			runCommand("sh", "-c", fmt.Sprintf("nmap -sC %s | sed 's/Nmap/NVSCANNER/g' | awk '{gsub(\"Nmap\",\"NVSCANNER\")}1' | sed 's/https:\\/\\/nmap.org//g'", target))
		case "6":
			fmt.Print("Enter port range (e.g., 1-1024): ")
			portRange, _ := reader.ReadString('\n')
			portRange = strings.TrimSpace(portRange)
			runCommand("sh", "-c", fmt.Sprintf("nmap -p %s %s | sed 's/Nmap/NVSCANNER/g' | sed '/https:\\/\\/nmap.org/d'", portRange, target))
		case "7":
			runCommand("sh", "-c", fmt.Sprintf("nmap -sV %s | sed 's/Nmap/NVSCANNER/g' | awk '{gsub(\"Nmap\",\"NVSCANNER\")}1' | sed '/https:\\/\\/nmap.org/d' | awk '/^PORT/{flag=1} flag; /^Nmap done/{flag=0}'", target))		
		case "8":
			runCommand("sh", "-c", fmt.Sprintf("sudo nmap -O --osscan-guess %s | sed 's/Nmap/NVSCANNER/g' | sed '/https:\\/\\/nmap.org/d' | awk '!/PORT|STATE|SERVICE/'", target))
		default:
			fmt.Println("Invalid option:", option)
		}
	}
	// Ask if the user wants to rerun the tool
	fmt.Print("Do you want to run the tool again? (y/n): ")
	rerunChoice, _ := reader.ReadString('\n')
	rerunChoice = strings.TrimSpace(rerunChoice)
	if rerunChoice == "y" || rerunChoice == "Y" {
		runTool()
	} else {
		fmt.Println("Exiting...")
	}
}


func scanPorts(target string) {
	fmt.Println("---------------------------------------------------------RESULT---------------------------------------------------------")
	// Define the range of ports to scan
	startPort := 1
	endPort := 1024

	fmt.Println("Scanning ports", startPort, "-", endPort)

	// Use a WaitGroup to ensure all port scans finish before returning
	var wg sync.WaitGroup
	wg.Add(endPort - startPort + 1)

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", target, port)
		go func(address string) {
			defer wg.Done()
			conn, err := net.Dial("tcp", address)
			if err == nil {

				fmt.Printf("Port %d: Open\n", port)
				conn.Close()
			}
		}(address)
	}

}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("success", output)
		return
	}
	fmt.Println("---------------------------------------------------------RESULT---------------------------------------------------------")
	fmt.Println(string(output))
}

