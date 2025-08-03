package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"time"
	"math/rand"

	"github.com/curator4/pokedexcli/internal/api"
	"github.com/curator4/pokedexcli/internal/pokecache"
)

type Config struct {
	LocationAreaNext string
	LocationAreaPrev string
	Cache pokecache.Cache
	Pokedex map[string]api.Pokemon
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

	cfg := &Config {
		Pokedex: make(map[string]api.Pokemon),
		Cache: pokecache.NewCache(10 * time.Second),
	}

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
			err := cliCommand.callback(cfg, args)
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
		"catch": {
			name: "catch",
			description: "Attempt to catch a a pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Inspect a caught pokemon",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "List caught pokemon",
			callback: commandPokedex,
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

func commandCatch(cfg *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("provide a pokemon name")
	}

	name := args[0]
	
	pokemon, err := api.GetPokemon(name, &cfg.Cache)
	if err != nil {
		fmt.Errorf("failed to get pokemon data: %w", err)
	}
	
	if pokemon == nil {
		return fmt.Errorf("pokemon %s not found\n", name)
	}

	catchRate := 255 / (1 + pokemon.Base_experience / 100)

	fmt.Print("\n")
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	if rand.Intn(256) < catchRate {
		fmt.Printf("%s was caugth!\n", name)
		fmt.Printf("You may now inspect it with the inspect command.\n\n")
		cfg.Pokedex[name] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n\n", name)
	}

	return nil
}

func commandInspect(cfg *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("provide a pokemon name")
	}

	name := args[0]

	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		fmt.Print("\nyou have not caught that pokemon\n\n")
		return nil
	}

	fmt.Printf("\n")
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, pokemonStat := range pokemon.Stats {
		fmt.Printf("    -%s: %d\n", pokemonStat.Stat.Name, pokemonStat.Base_stat)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("    - %s\n", pokemonType.Type.Name)
	}
	fmt.Printf("\n")

	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	fmt.Print("\n")
	fmt.Print("Your Pokedex:\n")
	for _, pokemon := range cfg.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	fmt.Print("\n")
	return nil
}
