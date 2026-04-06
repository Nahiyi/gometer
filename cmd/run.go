package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pressure test",
	RunE:  runPressureTest,
}

func runPressureTest(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}
