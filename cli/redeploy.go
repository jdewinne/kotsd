package cli

import (
	"fmt"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RedeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "redeploy [flags] [...name]",
		Short:         "Redeploy all application versions",
		Long:          `Redeploy all application versions`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configs := runtime_conf.Configs
			if len(args) > 0 {
				configs = filter(configs, args)
			}
			c := make(chan kotsd.Instance, len(configs))
			for _, instance := range configs {
				go redeployVersions(c, instance)
			}
			for range configs {
				i := <-c
				fmt.Println("Redeployed", i.Name)
			}

			return nil

		},
	}

	return cmd
}

func redeployVersions(c chan kotsd.Instance, instance kotsd.Instance) {
	err := instance.RedeployApps()
	if err != nil {
		fmt.Println("failed to redeploy instance", instance.Name, err)
	}
	c <- instance
}
