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
	callback    func([]string) error
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

type ExploreResults struct {
	Name               string `json:"name"`
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string
			Url  string
		}
	} `json:"pokemon_encounters"`
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
		"explore": {
			name:        "explore",
			description: "Show names of pokemon in location",
			callback:    commandExplore,
		},
	}
}

func commandHelp(parameters []string) error {
	fmt.Println("Type 'exit' to exit the program")

	return nil
}

func commandExit(parameters []string) error {
	fmt.Println("Closing the pokedex")
	os.Exit(0)

	return nil
}

func commandMap(parameters []string) error {
	mapLocations()
	locationURL = *locs.Next
	return nil
}

func commandMapB(parameters []string) error {
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

func commandExplore(parameters []string) error {
	fmt.Println("Exploring " + parameters[0])
	fmt.Println("Found pokemon:")
	url := locationURL + parameters[0]

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	results := ExploreResults{}
	err = json.Unmarshal(body, &results)

	if err != nil {
		return err
	}
	for _, pok := range results.Pokemon_encounters {
		fmt.Println("- " + pok.Pokemon.Name)
	}

	return nil
}
