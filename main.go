package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
