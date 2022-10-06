package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fabiofenoglio/excelconv/services"
)

var rootCmd = &cobra.Command{
	Use:   "conv",
	Short: "Convert excel files",
	Long:  `Convert excel files from an ugly format to a human readable one`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here
		err := services.Run("./input.xlsx")
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
