package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// ipdnsupdateCmd represents the ipdnsupdate command
var ipdnsupdateCmd = &cobra.Command{
	Use:   "ipdnsupdate {ipaddress}",
	Short: "Updates the Netbox ipaddress with the DNS name for the IP",
	Long: `Does a PTR lookup for the given IP address.  When found, the 
	address is retrieved from Netbox and the dns_name field populated.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := svc.IPdnsUpdate(args[0]); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ipdnsupdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipdnsupdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipdnsupdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
