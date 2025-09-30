package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wzhanjun/ges/internal/project"
)

var (
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	rootCmd = &cobra.Command{
		Use:     "ges",
		Short:   "GES: A skeleton for Echo",
		Long:    `GES: A skeleton for Echo`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
