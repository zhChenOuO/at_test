package main

import (
	cmd "amazing_talker/cmd"
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "server,scheduler"}

func main() {
	rootCmd.AddCommand(cmd.ServerCmd, cmd.SchedulerCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
