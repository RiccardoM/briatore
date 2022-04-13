package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	reportcmd "github.com/riccardom/briatore/cmd/report"

	junocmd "github.com/forbole/juno/v3/cmd"
)

func main() {
	// Setup logging to be textual
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Config the runner
	config := junocmd.NewConfig("briatore")

	// Build the root command
	rootCmd := &cobra.Command{
		Use:   "briatore",
		Short: "Briatore is a tax reporter helper for Cosmos-based chains",
	}
	rootCmd.AddCommand(
		reportcmd.GetReportCmd(),
	)

	exec := junocmd.PrepareRootCmd(config.GetName(), rootCmd)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
