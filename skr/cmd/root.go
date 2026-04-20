package cmd

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/philippeVV/skill-registry/skr/internal/config"
	"github.com/philippeVV/skill-registry/skr/internal/telemetry"
)

var (
	registryOverride  string
	verbose           bool
	cfg               *config.Config
	telemetryShutdown func(context.Context) error
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

		// Init telemetry
		shutdown, err := telemetry.Init(cmd.Context(), cfg.OTELEndpoint, verbose)
		if err != nil {
			log.Debug().Err(err).Msg("telemetry init failed")
		} else {
			telemetryShutdown = shutdown
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if telemetryShutdown != nil {
			ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Second)
			defer cancel()
			return telemetryShutdown(ctx)
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
