package cli

import (
	"fmt"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "update [flags] [...name]",
		Short:         "Update all application versions",
		Long:          `Update all application versions`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configs := runtime_conf.Configs
			slug, _ := cmd.Flags().GetString("slug")
			if len(args) > 0 {
				configs = filter(configs, args)
			}
			c := make(chan kotsd.Instance, len(configs))
			for _, instance := range configs {
				go updateVersions(c, instance, slug)
			}
			for range configs {
				i := <-c
				fmt.Println("Updated", i.Name)
			}

			return nil

		},
	}

	cmd.Flags().StringP("slug", "s", "", "Specify the application slug of the application you wish to update")

	return cmd
}

func updateVersions(c chan kotsd.Instance, instance kotsd.Instance, slug string) {
	err := instance.UpdateApps(slug)
	if err != nil {
		fmt.Printf("failed to update instance %s", instance.Name)
	}
	c <- instance
}
