package cli

import (
	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete-instance",
		Short:         "Delete kots instance",
		Long:          `Delete kots instance`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			runtime_conf.DeleteInstance(name)
			kotsd.WriteConfig(&runtime_conf, cfgFile)
			return nil

		},
	}
	cmd.Flags().StringP("name", "n", "", "Name of the kots instance (should be unique)")
	cmd.MarkFlagRequired("name")

	return cmd
}
