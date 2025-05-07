package cmd

import (
	"github.com/shayantrix/smartfuzz/pkg/controllers"
	"fmt"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use: "sql",
	Aliases: []string{"target"},
	Short: "Put an url in front of sql command 'example: smartfuzz sql 'test.com''", 
	Long: "At this stage you put a url, we give you the entry points which can be exploited by payloads",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string){
		fmt.Printf("Your URL is: %s", args[0])
		controllers.SqlInjection(args[0])
	},
}

func init(){
	rootCmd.AddCommand(urlCmd)
}

