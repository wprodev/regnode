package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "statu",
	Short: "Prints current host status in the cluster",
	Long: `Current host status in the cluster can return the more information about its roles, cluster, basic info about the machine, 
register command and token that was used to connect to cluster and more. Configure the output format by using -o|--output [yaml|json].`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status")
	},
}
