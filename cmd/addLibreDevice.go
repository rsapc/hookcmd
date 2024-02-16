/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// addLibreDeviceCmd represents the addLibreDevice command
var addLibreDeviceCmd = &cobra.Command{
	Use:   "addLibreDevice {ip address} {netbox model} {netbox model ID}",
	Short: "Adds the given device to LibreNMS",
	Long: `Adds the IP address given to LibreNMS.  The Netbox model and ID
	are then used to update the monitoring_id custom field and record 
	status in the Journal for the device / VM.
	`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		useHTML, _ := cmd.Flags().GetBool("html")
		if useHTML {
			startHTML("Adding %s:%s with IP %s to LibreNMS", args[1], args[2], args[0])
		}
		modelID, err := strconv.ParseInt(args[2], 0, 64)
		if err != nil {
			log.Fatalf("could not parse modelID: %v", err)
		}

		if err = svc.AddToLibreNMS(args[0], args[1], modelID); err != nil {
			if !useHTML {
				log.Fatal(err)
			}
		}
		if useHTML {
			endHTML()
		}
	},
}

func init() {
	rootCmd.AddCommand(addLibreDeviceCmd)
	addLibreDeviceCmd.Flags().BoolP("html", "x", false, "Return response as HTML")

}
