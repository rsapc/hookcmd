/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// libreMissingReportCmd represents the libreMissingReport command
var libreMissingReportCmd = &cobra.Command{
	Use:   "libreMissingReport",
	Short: "Generates a CSV of netbox devices that are not in LibreNMS",
	Long: `Returns all Netbox devices that are active, with a primary_ip 
	assigned, and do not have a monitoring_id set and the IP does not 
	exist in LibreNMS.  If the IP is found in LibreNMS the device_id is
	set as the monitoring_id in Netbox
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fname, err := cmd.Flags().GetString("output")
		if err != nil {
			fname = "-"
		}
		var out io.Writer
		if fname == "-" {
			out = os.Stdout
		} else {
			out, err = os.Create(fname)
			if err != nil {
				log.Fatalf("error opening %s: %v", fname, err)
			}
		}
		svc.MissingFromLibre(out)
	},
}

func init() {
	rootCmd.AddCommand(libreMissingReportCmd)
	libreMissingReportCmd.Flags().StringP("output", "o", "-", "Output filename")
}
