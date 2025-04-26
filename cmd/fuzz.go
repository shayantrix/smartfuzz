package cmd

import (
	"github.com/spf13/cobra"
)

var fuzzCmd = &cobra.Command{
	Use:	"fuzz",
	Aliases: []string{"fuzzing"},
	Short: "Test for fuzzing",
	Run:	func(cmd *cobra.Command, args []string){
		Fuzz()
	},
}

func init(){
	rootCmd.AddCommand(fuzzCmd)
}
