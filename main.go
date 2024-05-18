package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func()
}

func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandHelp() {
	fmt.Println("Type 'exit' to exit the program")
}

func commandExit() {
	fmt.Println("Closing the pokedex")
	os.Exit(0)
}

func printPrompt() {
	fmt.Printf("Pokedex > ")
}

func main() {
	commands := commands()
	reader := bufio.NewScanner(os.Stdin)

	printPrompt()
	for reader.Scan() {
		text := strings.ToLower(strings.TrimSpace(reader.Text()))
		if command, exists := commands[text]; exists {
			command.callback()
		} else {
			fmt.Println("Invalid command")
		}
		printPrompt()
	}
	fmt.Println()
}
