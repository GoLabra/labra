/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: StartServer,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartServer(cmd *cobra.Command, args []string) {

	for {
		fmt.Printf("@@@@@@@@@@@@@@@@@ ~ Start Server ~ @@@@@@@@@@@@@@@@@\n")

		run := exec.Command("go", "run", "main.go")
		run.Dir = "../app"
		stderr, _ := run.StderrPipe()

		err := run.Start()

		if err != nil {
			fmt.Println("Error starting server: ", err)
			return
		}

		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)

		go func() {
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
			}
		}()

		run.Wait()

		fmt.Printf("@@@@@@@@@@@@@@@@@ ~ Server Stopped ~ @@@@@@@@@@@@@@@@@\n\n\n\n")
	}
}
