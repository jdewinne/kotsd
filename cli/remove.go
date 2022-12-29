package cli

import (
	"fmt"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "remove [flags] [...name]",
		Short:         "Remove all the applications",
		Long:          `Remove all the applications`,
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
				go removeApps(c, instance)
			}
			for range configs {
				i := <-c
				fmt.Println("Removed apps from", i.Name)
			}

			return nil

		},
	}

	return cmd
}

func removeApps(c chan kotsd.Instance, instance kotsd.Instance) {
	err := instance.RemoveApps()
	if err != nil {
		fmt.Println("failed to remove apps", instance.Name, err)
	}
	c <- instance
}
