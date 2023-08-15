package cmd

import (
	"fmt"
	"log"

	commands "fourclover.org/cmd/cli/commands"
	common "fourclover.org/internal/common"
	config "fourclover.org/internal/config"
	logger "fourclover.org/internal/logger"
)

// Cmd is the main function for the CLI.
func Cmd(subCommand string) bool {
	/* Cmd is the main function for the CLI. It takes a subcommand as an argument and executes the corresponding function.*/

	var configDataInBytes []byte
	if config.YamlConfigurationStatus() { // Check if the YAML configuration file exists, if yes read it and convert it to bytes for later use
		log.Default().Println("INFO: FourClover: configuration file found")
		configData, _ := config.GetYamlConfigData()
		configDataInBytes, _ = common.GetConfigDataInBytes(configData)
	}

	if !config.GetFourCloverSuppressLogs() { // Start logger if not suppressed
		logger.NewLogger()
	}

	switch subCommand { // Check which subcommand was passed to the CLI and execute the corresponding function
	case "snapshot":
		return snapshot(configDataInBytes)
	case "compare":
		return compare(configDataInBytes)
	case "policy":
		return policy(configDataInBytes)
	case "version":
		return commands.Version()
	case "help":
		return commands.Help()
	case "demo":
		return commands.Demo()
	default:
		fmt.Printf("INFO: FourClover: missing subcommand (snapshot or compare or policy)\n      Run 'fourclover help' for usage.\n      Run 'fourclover demo' for demonstration.\n")
		return false
	}
}

func snapshot(configDataInBytes []byte) bool {
	status, err := commands.Snapshot(configDataInBytes)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
		return false
	}
	return status
}

func compare(configDataInBytes []byte) bool {
	status, err := commands.Compare(configDataInBytes)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
		return false
	}
	return status
}

func policy(configDataInBytes []byte) bool {
	status, err := commands.Policy(configDataInBytes)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
		return false
	}
	return status
}
