package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"time"

	"github.com/curator4/pokedexcli/internal/api"
	"github.com/curator4/pokedexcli/internal/pokecache"
)

type Config struct {
	LocationAreaNext string
	LocationAreaPrev string
	Cache pokecache.Cache
}

type cliCommand struct {
	name string
	description string
	callback func(*Config, []string) error
}

var commands map[string]cliCommand
var cfg Config


func main () {
	initCommands()
	cfg.Cache = pokecache.NewCache(10 * time.Second)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)

		command := input[0]
		var args []string
		if len(input) > 1 {
			args = input[1:]
		}
		
		if cliCommand, ok := commands[command]; ok {
			err := cliCommand.callback(&cfg, args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
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
		"explore": {
			name: "explore",
			description: "Shows pokemon in area",
			callback: commandExplore,
		},
	}
}

func commandExit(cfg *Config, args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args []string) error {
	fmt.Print("Welcome to the Pokedex!\n\n")
	fmt.Print("Usage:\n")
	
	for _, cliCommand := range commands {
		fmt.Printf("%s: %s\n", cliCommand.name, cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}

func commandMap(cfg *Config, args []string) error {
	if cfg.LocationAreaNext == "" {
		cfg.LocationAreaNext = "https://pokeapi.co/api/v2/location-area/?limit=20" 
	}

	page, err := api.GetPage[api.LocationArea](cfg.LocationAreaNext, &cfg.Cache)
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

func commandMapB(cfg *Config, args []string) error {
	if cfg.LocationAreaPrev == "" {
		fmt.Printf("\nyou're on the first page\n")
		return nil
	}

	page, err := api.GetPage[api.LocationArea](cfg.LocationAreaPrev, &cfg.Cache)
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

func commandExplore(cfg *Config, args []string) error {
	area := args[0]
	areaData, err := api.GetAreaPokemon(area, &cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to get area data: %w", err)
	}


	fmt.Print("\n")
	fmt.Printf("Exploring %s\n", area)
	fmt.Print("Found Pokemon:\n")
	for _, encounter := range areaData.Pokemon_encounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	fmt.Print("\n")
	return nil
}
