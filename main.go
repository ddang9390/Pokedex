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
		t := strings.Fields(text)

		if command, exists := commands[t[0]]; exists {
			if len(t) > 1 {
				command.callback(t[1:])
			} else {
				command.callback(t)
			}

		} else {
			fmt.Println("Invalid command")
		}
		printPrompt()
	}
	fmt.Println()
}
