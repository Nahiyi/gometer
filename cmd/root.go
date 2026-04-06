package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	threads        int
	rampUp         int
	loop           int
	configFile     string
	outputFile     string
	requestTimeout int
	maxDuration    int
	dryRun         bool
)

var rootCmd = &cobra.Command{
	Use:   "gometer",
	Short: "gometer is a HTTP pressure testing tool",
	Long:  `gometer is a CLI HTTP pressure testing tool for learning Go standard library.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&threads, "threads", "n", 0, "Number of threads (required)")
	runCmd.Flags().IntVarP(&rampUp, "ramp-up", "t", 0, "Ramp-up period in seconds (default 0)")
	runCmd.Flags().IntVarP(&loop, "loop", "l", 1, "Number of loops per thread (default 1)")
	runCmd.Flags().StringVarP(&configFile, "config", "c", "./req.json", "Request config file (default ./req.json)")
	runCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default stdout)")
	runCmd.Flags().IntVar(&requestTimeout, "request-timeout", 5000, "Request timeout in milliseconds (default 5000)")
	runCmd.Flags().IntVar(&maxDuration, "max-duration", 0, "Max duration in seconds, 0 means unlimited (default 0)")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate config file without running")

	rootCmd.MarkFlagRequired("threads")
}
