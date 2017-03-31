package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	sfballetRoot = "https://www.sfballet.org"
)

var (
	programPaths []string
	programsMap  map[string][]*Program
)

func main() {
	doc, err := requestToNode(getRequest(sfballetRoot))
	if err != nil {
		log.Fatal(err)
	}
	programPaths = fetchTicketPaths(doc)

	var wg sync.WaitGroup
	channelsMap := make(map[string](chan []*Program))
	for _, path := range programPaths {
		channelsMap[path] = make(chan []*Program)
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			doc, err := requestToNode(getRequest(fmt.Sprintf("%s%s", sfballetRoot, path)))
			if err != nil {
				log.Fatal(err)
			}
			channelsMap[path] <- fetchPrograms(doc)
		}(path)
	}

	programsMap = make(map[string][]*Program)
	for path, c := range channelsMap {
		go func(path string, c chan []*Program) {
			if programs := <-c; len(programs) > 0 {
				for _, program := range programs {
					wg.Add(1)
					go func(program *Program) {
						defer wg.Done()
						doc, err := requestToNode(getRequest(fmt.Sprintf("%s", program.TicketURL)))
						if err != nil {
							log.Fatal(err)
						}
						program.ShowDates = fetchShowDates(doc)
					}(program)
				}
				programsMap[path] = programs
			}
		}(path, c)
	}
	fmt.Println("Loading...")
	wg.Wait()
	currState := NewState()
	fmt.Println("Simple sf ballet scraper")
	fmt.Println("------------------------")
	help(currState)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	for scanner.Scan() {
		text := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if text == "state" {
			fmt.Println(currState)
		} else if text == "help" {
			help(currState)
		} else if currState.IsBase() && text == "list paths" {
			for id, path := range programPaths {
				fmt.Printf("\t%d. %s\n", id, path)
			}
		} else if currState.IsBase() && strings.HasPrefix(text, "go path ") {
			pathIDStr := strings.TrimLeft(text, "go path ")
			pathID, err := strconv.Atoi(pathIDStr)
			if err != nil || (pathID >= len(programPaths) || pathID < 0) {
				fmt.Printf("Invalid path_id: %s\n", pathIDStr)
				fmt.Print(">>> ")
				continue
			}
			currState.SetPath(programPaths[pathID])
			help(currState)
		} else if currState.IsPath() && text == "list programs" {
			for id, program := range programsMap[currState.Path] {
				fmt.Printf("\t%d: %s", id, program.Title)
				if !program.Available {
					fmt.Print(" (not currently available)")
				} else {
					fmt.Printf(" %s", program.Dates)
				}
				fmt.Println()
			}
		} else if currState.IsPath() && strings.HasPrefix(text, "go program ") {
			programIDStr := strings.TrimLeft(text, "go program ")
			programID, err := strconv.Atoi(programIDStr)
			if err != nil || (programID >= len(programsMap[currState.Path]) || programID < 0) {
				fmt.Printf("Invalid program_id: %s\n", programIDStr)
				fmt.Print(">>> ")
				continue
			}
			currState.SetProgram(programID)
			help(currState)
		} else if currState.IsProgram() && text == "list dates" {
			for _, date := range programsMap[currState.Path][currState.Program].ShowDates {
				fmt.Printf("\t%v\n", date)
			}
		} else if (currState.IsPath() || currState.IsProgram()) && text == "back" {
			switch currState.Type {
			case pathType:
				currState.UnsetPath()
				fmt.Println("Going back to choose a path.")
			case programType:
				currState.UnsetProgram()
				fmt.Println("Going back to choose a program.")
			}
			help(currState)
		} else if text != "" {
			fmt.Println("Unrecognized command. Try 'help' for a list of available commands.")
		}
		fmt.Print(">>> ")
	}
	fmt.Println()
	fmt.Println("Good bye!")
}

func help(state *State) {
	fmt.Println("Here are the available commands:")
	fmt.Println("\thelp")
	if state.Base && state.IsBase() {
		fmt.Println("\tlist paths")
		fmt.Println("\tgo path <path_id>, where path_id is the number for one of the following:")
		for id, path := range programPaths {
			fmt.Printf("\t\t%d. %s\n", id, path)
		}
	} else if state.IsPath() {
		fmt.Println("\tback")
		if _, ok := programsMap[state.Path]; ok {
			fmt.Println("\tlist programs")
			fmt.Println("\tgo program <program_id>, where program_id is the number for one of the following:")
			for id, program := range programsMap[state.Path] {
				fmt.Printf("\t\t%d: %ss", id, program.Title)
				if !program.Available {
					fmt.Print(" (not currently available)")
				} else {
					fmt.Printf(" %s", program.Dates)
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("Note: No programs available.\n")
		}
	} else if state.IsProgram() {
		fmt.Println("\tback")
		if _, ok := programsMap[state.Path]; ok {
			fmt.Println("\tlist dates")
		} else {
			fmt.Println("Note: No dates available.\n")
		}
	}
}
