package cmd

import (
        "github.com/shayantrix/smartfuzz/pkg/controllers"
        "fmt"
        "github.com/spf13/cobra"
)

var commandCmd = &cobra.Command{
        Use: "commandInjection",
        Aliases: []string{"target", "command", "injection"},
        Short: "Put an url in front of commandInjection 'example: smartfuzz commandInjection 'test.com''",
        Long: "You will fuzz the target with command injection payloads. The payload may vary due to the target's server",
        Args: cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string){
                fmt.Printf("Your URL is: %s", args[0])
                controllers.CommandInjection(args[0])
        },
}

func init(){
        rootCmd.AddCommand(commandCmd)
}

