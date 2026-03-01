package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sheldon <command>")
		return
	}

	switch os.Args[1] {
	case "record":
		handleRecordCommand()
	case "search":
		handleSearchCommand()
	default:
		fmt.Println("Unknown command.")
	}
}

func handleRecordCommand() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: sheldon record <command> <cwd> <exitcode>")
		return
	}

	command := os.Args[2]
	cwd := os.Args[3]

	exitcode, err := strconv.Atoi(os.Args[4])
	if err != nil {
		// TODO: Build a logging system.
		fmt.Println("Error: Invalid exit code")
		return
	}

	err = RecordCommand(command, cwd, exitcode)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func handleSearchCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sheldon search <query>")
		return
	}

	query := os.Args[2]

	results, err := SearchCommand(query)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, command := range results {
		fmt.Println(command)
	}
}
