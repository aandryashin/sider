package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	siderURL string
	timeout  time.Duration
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&siderURL, "url", "u", "http://localhost:8080", "Sider API URL")
	RootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "Operation timeout")
	RootCmd.AddCommand(keysCmd)
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(setCmd)
	RootCmd.AddCommand(delCmd)
}

var RootCmd = &cobra.Command{
	Use:   "sider",
	Short: "Client for Sider API",
	Long:  "Client for Redis like key value in-memory storage.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
