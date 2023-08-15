package main

import (
	"log"

	cli "fourclover.org/cmd/cli"
	config "fourclover.org/internal/config"
)

var subCommand string

func main() {
	if config.GetFourCloverIsSnapshot() {
		subCommand = "snapshot"
	} else if config.GetFourCloverIsCompare() {
		subCommand = "compare"
	} else if config.GetFourCloverIsPolicy() {
		subCommand = "policy"
	} else if config.GetFourCloverVersionArg() {
		subCommand = "version"
	} else if config.GetFourCloverHelp() {
		subCommand = "help"
	} else if config.GetFourCloverDemo() {
		subCommand = "demo"
	} else {
		log.Fatalf("ERROR: FourClover: missing subcommand (snapshot or compare or policy)\n      Run 'fourclover help' for usage.\n      Run 'fourclover demo' for demonstration.\n")
	}

	if config.GetFourCloverIsSnapshot() && config.GetFourCloverIsCompare() {
		log.Fatalf("ERROR: FourClover: cannot run both snapshot and compare at the same time")
	}

	cli.Cmd(subCommand) // Call the main function for the CLI. It takes a subcommand as an argument and executes the corresponding function.
}
