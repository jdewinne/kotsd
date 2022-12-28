package cli

import (
	"fmt"
	"os"
	"strings"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list [flags] [...name]",
		Short:         "List all kots instance version and application versions",
		Long:          `List all kots instance version and application versions`,
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
			t.AppendHeader(table.Row{"Name", "Kots Version", "Connection", "#", "Application Name", "Version", "Upgrades"})
			for range configs {
				i := <-c
				if len(i.Apps) == 0 {
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error},
					})
				}
				for indx, app := range i.Apps {
					pversions := "-"
					if len(app.PendingVersions) > 0 {
						pversionstrings := []string{}
						for _, pversion := range app.PendingVersions {
							pversionstrings = append(pversionstrings, fmt.Sprintf("%d - %s", pversion.Sequence, pversion.VersionLabel))
						}
						pversions = strings.Join(pversionstrings, "\n")
					}
					t.AppendRows([]table.Row{
						{i.Name, i.KotsVersion, i.Error, indx, app.Name, app.Version, pversions},
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
			versionLabel := "Not defined"
			if app.Downstream.CurrentVersion != nil {
				versionLabel = app.Downstream.CurrentVersion.VersionLabel
			}
			application := kotsd.Application{Name: app.Name, Version: versionLabel}
			var pVersions []kotsd.PendingVersion
			for _, pv := range app.Downstream.PendingVersions {
				pVersions = append(pVersions, kotsd.PendingVersion{VersionLabel: pv.VersionLabel, Sequence: pv.Sequence})
			}
			application.PendingVersions = pVersions
			i.Apps = append(i.Apps, application)
		}
	}

	c <- i
}
