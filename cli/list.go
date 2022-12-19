package cli

import (
	"os"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all kots instance version and application versions",
		Long:          `List all kots instance version and application versions`,
		ArgAliases:    []string{"name"},
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
				go getVersions(c, instance)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleColoredBlackOnBlueWhite)
			t.AppendHeader(table.Row{"Name", "Kots Version", "Connection", "#", "Application Name", "Version", "Upgrade Available"})
			for range configs {
				i := <-c
				if len(i.Apps) == 0 {
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error},
					})
				}
				for indx, app := range i.Apps {
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error, indx, app.Name, app.Version, app.PendingVersion},
					})
				}
				t.AppendSeparator()
			}
			t.Render()
			return nil

		},
	}

	return cmd
}

func filter(instances []kotsd.Instance, args []string) []kotsd.Instance {
	var configs []kotsd.Instance
	for _, i := range instances {
		for _, name := range args {
			if i.Name == name {
				configs = append(configs, i)
				break
			}
		}
	}
	return configs
}

func getVersions(c chan kotsd.Instance, i kotsd.Instance) {
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
