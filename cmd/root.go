package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "smartfuzz",
    Short: "smartfuzz is a cli tool",
    Run: func(cmd *cobra.Command, args []string) {

    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
        os.Exit(1)
    }
}
