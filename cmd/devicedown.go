package cmd

import (
	"github.com/spf13/cobra"
)

// devicedownCmd represents the devicedown command
var devicedownCmd = &cobra.Command{
	Use:   "devicedown {alert payload}",
	Short: "Sets the Netbox status to offline when it goes down in LibreNMS",
	Long: `Receives an alert from LibreNMS and sets the corresponding device 
	in Netbox to Offline`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc.DeviceDown(args[0])
	},
}

func init() {
	rootCmd.AddCommand(devicedownCmd)

}
