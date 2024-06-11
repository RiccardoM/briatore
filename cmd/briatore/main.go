package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	reportcmd "github.com/riccardom/briatore/cmd/report"
	startcmd "github.com/riccardom/briatore/cmd/start"
	"github.com/riccardom/briatore/utils"
)

func main() {
	// Setup logging to be textual
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Build the root command
	rootCmd := &cobra.Command{
		Use:   "briatore",
		Short: "Briatore is a tax reporter helper for Cosmos-based chains",
	}
	rootCmd.AddCommand(
		reportcmd.GetReportCmd(),
		startcmd.GetStartCmd(),
	)

	exec := prepareRootCmd("briatore", rootCmd)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// PrepareRootCmd is meant to prepare the given command binding all the viper flags
func prepareRootCmd(name string, cmd *cobra.Command) cli.Executor {
	cmd.PersistentPreRunE = utils.ConcatCobraCmdFuncs(
		utils.BindFlagsLoadViper,
		cmd.PersistentPreRunE,
	)

	home, _ := os.UserHomeDir()
	defaultConfigPath := path.Join(home, fmt.Sprintf(".%s", name))
	cmd.PersistentFlags().String("home", defaultConfigPath, "Set the home folder of the application, where all files will be stored")

	return cli.Executor{Command: cmd, Exit: os.Exit}
}
