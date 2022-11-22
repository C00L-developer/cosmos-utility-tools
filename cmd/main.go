package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	// app name
	appName = "cosmos-utility"
	// version represents the program based on the git tag
	version = "v0.1.0"

	// flags
	// chain name flag
	flagChain = "chainName"
	// account prefix flag
	flagPrefix = "accountPrefix"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version

	monitorFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     flagChain,
			Aliases:  []string{"c"},
			Usage:    "Cosmos based chain name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagPrefix,
			Aliases:  []string{"p"},
			Usage:    "Account prefix name",
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "validator-status",
			Aliases: []string{},
			Usage:   "Run the monitoring tool of validators status",
			Action:  monitorValidator,
			Flags:   monitorFlags,
		},
		{
			Name:    "vesting-analyze",
			Aliases: []string{},
			Usage:   "Analyze vesting accounts",
			Action:  vestingAnalyze,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}
