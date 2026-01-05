package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func handleConfig(args []string) {
	if len(args) == 0 {
		showConfigUsage()
		return
	}

	switch args[0] {
	case "set-key":
		handleSetKey()
	case "remove-key":
		handleRemoveKey()
	case "status":
		handleStatus()
	default:
		showConfigUsage()
	}
}

func handleSetKey() {
	fmt.Print("Enter your Gemini API key: ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error reading key: %v\n", err)
		return
	}
	fmt.Println()

	key := strings.TrimSpace(string(password))

	if key == "" {
		fmt.Println("API key cannot be empty")
		return
	}

	fmt.Print("Confirm your API key (press Enter to skip confirmation): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "" && response != "y" && response != "yes" {
		fmt.Print("Re-enter your API key: ")
		confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("Error reading key: %v\n", err)
			return
		}
		fmt.Println()

		if string(confirmPassword) != key {
			fmt.Println("Keys do not match")
			return
		}
	}

	err = SaveAPIKey(key)
	if err != nil {
		fmt.Println("⚠️ Could not save API key securely (keyring unavailable)")
		return
	}

	fmt.Println("API key saved securely")
}

func handleRemoveKey() {
	if !HasAPIKey() {
		fmt.Println("ℹ️ No API key configured")
		return
	}

	fmt.Print("Are you sure you want to remove your API key? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("Cancelled")
		return
	}

	err := DeleteAPIKey()
	if err != nil {
		return
	}

	fmt.Println("API key removed")
}

func handleStatus() {
	if os.Getenv("GEMINI_API_KEY") != "" {
		fmt.Println("Gemini API key is configured (environment variable)")
		return
	}

	if HasAPIKey() {
		fmt.Println("Gemini API key is configured (keyring)")
		return
	}

	fmt.Println("✗ No Gemini API key configured")
}

func showConfigUsage() {
	fmt.Println("usage: ship config <subcommand>")
	fmt.Println("")
	fmt.Println("Subcommands:")
	fmt.Println("  set-key       Securely store your Gemini API key")
	fmt.Println("  remove-key    Remove stored API key")
	fmt.Println("  status        Check if API key is configured")
}
