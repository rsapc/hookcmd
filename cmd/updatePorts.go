/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// updatePortsCmd represents the updatePorts command
var updatePortsCmd = &cobra.Command{
	Use:   "updatePorts {netbox ID} {libreNMS ID}",
	Short: "Updates the Netbox port descriptions from LibreNMS",
	Long: `Updates the interface descriptions on the Netbox device ID
	to match what is in LibreNMS device
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		nbID, err := strconv.ParseInt(args[0], 0, 0)
		if err != nil {
			log.Fatalf("could not parse netbox device ID: %v", err)
		}
		libreID, err := strconv.ParseInt(args[1], 0, 0)
		if err != nil {
			log.Fatalf("could not parse monitoring ID: %v", err)
		}
		useHTML, _ := cmd.Flags().GetBool("html")
		if useHTML {
			startHTML("Updating from LibreNMS device %v", nbID)
		}
		svc.UpdatePortDescriptions(int(nbID), int(libreID))
		if useHTML {
			endHTML()
		}
	},
}

func init() {
	rootCmd.AddCommand(updatePortsCmd)
	updatePortsCmd.Flags().BoolP("html", "x", false, "Return response as HTML")
}
