package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "add-instance",
		Short:         "Add kots instance",
		Long:          `Add kots instance`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil

		},
	}
	cmd.Flags().StringP("name", "n", "", "Name of the kots instance (should be unique)")
	cmd.Flags().StringP("endpoint", "e", "", "URL of the kots instance, for example http://10.10.10.5:8800")
	cmd.Flags().StringP("password", "p", "", "Password of the kots instance")
	cmd.Flags().BoolP("tlsVerify", "v", true, "If false, insecure or self signed tls for the kots instance will be allowed")
	return cmd
}
