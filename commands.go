package main

import (
	"Pokedex/internal/pokecache"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var locationURL string = "https://pokeapi.co/api/v2/location-area/"
var locs = Locations{}
var cache = pokecache.NewCache(time.Duration(5))

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
		"map": {
			name:        "map",
			description: "Showing names of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go backwards for locations list",
			callback:    commandMapB,
		},
	}
}

func commandHelp() error {
	fmt.Println("Type 'exit' to exit the program")

	return nil
}

func commandExit() error {
	fmt.Println("Closing the pokedex")
	os.Exit(0)

	return nil
}

func commandMap() error {
	mapLocations()
	locationURL = *locs.Next
	return nil
}

func commandMapB() error {
	if locs.Previous == nil {
		fmt.Println("You are on the first page")
		return errors.New("on the first page")
	}
	locationURL = *locs.Previous
	mapLocations()
	return nil
}

func mapLocations() error {
	body, ok := cache.Get(locationURL)
	if !ok {
		res, err := http.Get(locationURL)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		res.Body.Close()
		cache.Add(locationURL, body)
	}

	err := json.Unmarshal(body, &locs)
	if err != nil {
		return err
	}

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
