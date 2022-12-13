package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list-instances",
		Short:         "List all kots instances",
		Long:          `List all kots instances`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Instances:")
			for _, instance := range runtime_conf.Configs {
				fmt.Println("Name:", instance.Name, "- Endpoint:", instance.Endpoint)
			}
			return nil

		},
	}

	return cmd
}
