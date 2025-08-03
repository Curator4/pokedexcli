package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"

	"github.com/curator4/pokedexcli/internal/api"
)

type Config struct {
	LocationAreaNext string
	LocationAreaPrev string
}

type cliCommand struct {
	name string
	description string
	callback func(*Config) error
}

var commands map[string]cliCommand
var cfg Config


func main () {
	initCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		command := input[0]
		

		
		cliCommand, ok := commands[command]
		if ok {
			err := cliCommand.callback(&cfg)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		} else {
			fmt.Print("Unknown command\n\n")
		}
	}
}



func cleanInput(text string) []string {
	trimmed := strings.Trim(text, " ")
	lowered := strings.ToLower(trimmed)
	fields := strings.Fields(lowered)

	return fields
}


func initCommands() {
	commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Displays next page of location areas",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays previous page of location areas",
			callback: commandMapB,
		},
	}
}

func commandExit(cfg *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Print("Welcome to the Pokedex!\n\n")
	fmt.Print("Usage:\n")
	
	for _, cliCommand := range commands {
		fmt.Printf("%s: %s\n", cliCommand.name, cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}

func commandMap(cfg *Config) error {
	if cfg.LocationAreaNext == "" {
		cfg.LocationAreaNext = "https://pokeapi.co/api/v2/location-area/?limit=20" 
	}

	page, err := api.GetPage[api.LocationArea](cfg.LocationAreaNext)
	if err != nil {
		fmt.Printf("Error getting Location Areas: %v", err)
	}
	cfg.LocationAreaNext = page.Next
	cfg.LocationAreaPrev = page.Previous

	fmt.Printf("\n")
	for _, locationArea := range page.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}
	fmt.Printf("\n")

	return nil
}

func commandMapB(cfg *Config) error {
	if cfg.LocationAreaPrev == "" {
		fmt.Printf("\nyou're on the first page\n")
		return nil
	}

	page, err := api.GetPage[api.LocationArea](cfg.LocationAreaPrev)
	if err != nil {
		fmt.Printf("Error getting Location Areas: %v", err)
	}
	cfg.LocationAreaNext = page.Next
	cfg.LocationAreaPrev = page.Previous

	fmt.Printf("\n")
	for _, locationArea := range page.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}
	fmt.Printf("\n")

	return nil
}
