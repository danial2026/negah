package main

import (
	"bufio"
	"fmt"
	"negah/scanner"
	"os"
	"strconv"
	"strings"
)

func main() {
	if !scanner.CheckNmap() {
		fmt.Println("Heads up: 'nmap' isn't in your PATH.")
		fmt.Println("You'll need it for most of these scans. (brew install nmap / pacman -S nmap)")
	}

	reader := bufio.NewReader(os.Stdin)
	tools := scanner.GetFeatures()

	for {
		fmt.Println("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("                    NEGAH                    ")
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

		for _, t := range tools {
			dots := strings.Repeat(" ", 25-len(t.Name))
			fmt.Printf("%2d. %s %s %s\n", t.ID, t.Name, dots, t.Description)
		}

		fmt.Println(" 0. Exit")
		fmt.Print("\nWhat is your command? ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			fmt.Println("Goodbye.")
			break
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(tools) {
			fmt.Println("That wasn't a valid option. Try one of the numbers above.")
			continue
		}

		selected := tools[choice-1]
		target := ""

		// Some tools are self-contained and don't need a target domain/IP.
		if selected.ID != 12 && selected.ID != 13 {
			fmt.Print("Provide a target (IP, Domain, or Range): ")
			target, _ = reader.ReadString('\n')
			target = strings.TrimSpace(target)
		}

		// If they picked the custom range tool, we need to ask what ports they want.
		if selected.ID == 4 {
			fmt.Print("Which ports? (like 80,443 or 1-1000): ")
			ports, _ := reader.ReadString('\n')
			ports = strings.TrimSpace(ports)
			selected.Command = "-p " + ports
		}

		scanner.RunScan(selected, target)

		fmt.Print("\nDone. Press Enter to return to Negah...")
		reader.ReadString('\n')
	}
}
