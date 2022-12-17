package cli

import (
	"fmt"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all kots instance version and application versions",
		Long:          `List all kots instance version and application versions`,
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Instances:")
			c := make(chan kotsd.Instance, len(runtime_conf.Configs))
			for _, instance := range runtime_conf.Configs {
				go getKotsVersion(c, instance)
			}
			for range runtime_conf.Configs {
				i := <-c
				fmt.Println("Name:", i.Name, "- Kots version:", i.KotsVersion, "Connection:", i.Error)
				for _, app := range i.Apps {
					fmt.Println("Application:", app.Name, "- Version:", app.Version, "- Upgrade Available:", app.PendingVersion)
				}
			}
			return nil

		},
	}

	return cmd
}

func getKotsVersion(c chan kotsd.Instance, i kotsd.Instance) {
	kh, err := i.GetKotsHealthz()
	if err != nil {
		i.Error = err.Error()
	} else {
		i.KotsVersion = kh.Version
	}
	apps, err := i.GetApps()
	if err != nil {
		i.Error = err.Error()
	} else {
		for _, app := range apps.Apps {
			application := kotsd.Application{Name: app.Name, Version: app.Downstream.CurrentVersion.VersionLabel}
			var pVersions []string
			for _, pv := range app.Downstream.PendingVersions {
				pVersions = append(pVersions, pv.VersionLabel)
			}
			application.PendingVersion = pVersions
			i.Apps = append(i.Apps, application)
		}
	}

	c <- i
}
