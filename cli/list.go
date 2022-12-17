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
		SilenceUsage:  true,
		SilenceErrors: false,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := make(chan kotsd.Instance, len(runtime_conf.Configs))
			for _, instance := range runtime_conf.Configs {
				go getVersions(c, instance)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleColoredBlackOnBlueWhite)
			t.AppendHeader(table.Row{"Name", "Kots Version", "Connection", "Application Name", "Version", "Upgrade Available"})
			for range runtime_conf.Configs {
				i := <-c
				if len(i.Apps) == 0 {
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error},
					})
				}
				for _, app := range i.Apps {
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error, app.Name, app.Version, app.PendingVersion},
					})
				}
				t.AppendSeparator()
			}
			t.AppendSeparator()
			t.Render()
			return nil

		},
	}

	return cmd
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
