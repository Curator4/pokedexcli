package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands map[string]cliCommand


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
			err := cliCommand.callback()
			if err != nil {
				fmt.Printf("error: %w", err)
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
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n\n")
	fmt.Print("Usage:\n")
	
	for _, cliCommand := range commands {
		fmt.Printf("%s: %s\n", cliCommand.name, cliCommand.description)
	}
	fmt.Print("\n")
	return nil
}
