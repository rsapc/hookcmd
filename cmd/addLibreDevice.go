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
		modelID, err := strconv.ParseInt(args[2], 0, 64)
		if err != nil {
			log.Fatalf("could not parse modelID: %v", err)
		}

		if err = svc.AddToLibreNMS(args[0], args[1], modelID); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addLibreDeviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addLibreDeviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addLibreDeviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
