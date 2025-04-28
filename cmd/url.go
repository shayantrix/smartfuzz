package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use: "url",
	Aliases: []string{"target"},
	Short: "url to fuzz", 
	Long: "At this stage you put a url, we give you the entry points which can be exploited by payloads",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string){
		fmt.Printf("Your URL is: %s", args[0])
		controllers.URLHandler(args[0])
	},
}

func init(){
	rootCmd.AddCommand(urlCmd)
}

