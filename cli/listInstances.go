package cli

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleColoredBlackOnBlueWhite)
			t.AppendHeader(table.Row{"Name", "Endpoint"})
			for _, instance := range runtime_conf.Configs {
				t.AppendRows([]table.Row{{instance.Name, instance.Endpoint}})
			}
			t.Render()
			return nil

		},
	}

	return cmd
}
