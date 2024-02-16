/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// updatedeviceCmd represents the updatedevice command
var updatedeviceCmd = &cobra.Command{
	Use:   "updatedevice {monitoring_id}",
	Short: "Updates Netbox for the given LibreNMS ID",
	Long: `Looks up the provided device_id in LibreNMS.  If found 
	it will find the corresponding device in Netbox and update it 
	with discovered information from LibreNMS
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deviceID, err := strconv.ParseInt(args[0], 0, 0)
		if err != nil {
			log.Fatalf("could not parse device ID: %v", err)
		}
		useHTML, _ := cmd.Flags().GetBool("html")
		if useHTML {
			startHTML("Updating from LibreNMS device %v", deviceID)
		}
		svc.GetDeviceInfo(int(deviceID))
		if useHTML {
			endHTML()
		}
	},
}

func init() {
	rootCmd.AddCommand(updatedeviceCmd)
	updatedeviceCmd.Flags().BoolP("html", "x", false, "Return response as HTML")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updatedeviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updatedeviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
