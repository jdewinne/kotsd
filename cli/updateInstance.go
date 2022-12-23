package cli

import (
	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UpdateInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "update-instance",
		Short:         "Update kots instance",
		Long:          `Update kots instance`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			instance, err := runtime_conf.GetInstance(name)
			if err != nil {
				return errors.Wrap(err, "Did not find instance")
			}
			endpoint, _ := cmd.Flags().GetString("endpoint")
			if !cmd.Flags().Changed("endpoint") {
				endpoint, _ = PromptForEndpoint(instance.Endpoint)
			}
			tlsVerify, _ := cmd.Flags().GetBool("tlsVerify")
			if !cmd.Flags().Changed("tlsVerify") {
				tlsVerify, _ = PromptForTlsVerify(!instance.InsecureSkipVerify)
			}
			password, _ := PromptForPassword()
			err = runtime_conf.UpdateInstance(name, endpoint, password, tlsVerify)
			if err != nil {
				return err
			}
			kotsd.WriteConfig(&runtime_conf, cfgFile)
			return nil

		},
	}
	cmd.Flags().StringP("name", "n", "", "Name of the kots instance (should be unique)")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringP("endpoint", "e", "", "URL of the kots instance, for example http://10.10.10.5:8800")

	cmd.Flags().BoolP("tlsVerify", "v", true, "If false, insecure or self signed tls for the kots instance will be allowed")

	return cmd
}
