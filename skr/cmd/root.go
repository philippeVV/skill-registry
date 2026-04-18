package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/config"
)

var (
	registryOverride string
	verbose          bool
	cfg              *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "skr",
	Short: "Skill registry CLI — install, manage, and discover Claude Code packages",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Logging
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		// Load config
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Apply registry override
		if registryOverride != "" {
			cfg.Registry = registryOverride
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&registryOverride, "registry", "", "override registry URL for this invocation")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
