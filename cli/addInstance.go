package cli

import (
	kotsd "github.com/jdewinne/kotsd/pkg"
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
			name, _ := cmd.Flags().GetString("name")
			endpoint, _ := cmd.Flags().GetString("endpoint")
			tlsVerify, _ := cmd.Flags().GetBool("tlsVerify")
			password, _ := PromptForNewPassword()
			runtime_conf.AddInstance(name, endpoint, password, tlsVerify)
			kotsd.WriteConfig(&runtime_conf, cfgFile)
			return nil

		},
	}
	cmd.Flags().StringP("name", "n", "", "Name of the kots instance (should be unique)")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringP("endpoint", "e", "", "URL of the kots instance, for example http://10.10.10.5:8800")
	cmd.MarkFlagRequired("endpoint")

	cmd.Flags().BoolP("tlsVerify", "v", true, "If false, insecure or self signed tls for the kots instance will be allowed")

	return cmd
}
