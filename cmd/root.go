package cmd

import (
	"os"

	"github.com/rsapc/hookcmd/service"
	"github.com/spf13/cobra"
)

var svc *service.Service

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Various webhooks for interfacing with Netbox and LibreNMS",
	Long: `This program is designed to be called from a webhook service.
	It has been specifically designed to work with https://github.com/adnanh/webhook.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webhooks.yaml)")

	svc = service.NewService(os.Getenv, nil)
}
