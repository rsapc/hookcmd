/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// updatebyipCmd represents the updatebyip command
var updatebyipCmd = &cobra.Command{
	Use:   "updatebyip {ip}",
	Short: "Finds the IP in LibreNMS and updates the corresponding device in Netbox",
	Long: `Looks up the IP in Netbox to ensure it has an assigned object.  Then the 
	IP is searched for in LibreNMS.  If it is found the device is updated in Netbox
	with some of the discovered information in LibreNMS.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useHTML, _ := cmd.Flags().GetBool("html")
		if useHTML {
			startHTML("Updating from LibreNMS IP %s", args[0])
		}
		svc.FindDevice(args[0])
		if useHTML {
			endHTML()
		}
	},
}

func init() {
	rootCmd.AddCommand(updatebyipCmd)
	updatebyipCmd.Flags().BoolP("html", "x", false, "Return response as HTML")
}
