package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bootable",
	Short: "A tool to create multiboot USB drives",
	Long: `bootable is a powerful and easy-to-use command-line
	utility for creating bootable USB drives with multiple operating systems.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
