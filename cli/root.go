package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	kotsd "github.com/jdewinne/kotsd/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var runtime_conf kotsd.KotsdConfig
var version = "0.0.1"

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kotsd",
		Version:      version,
		Short:        "Run commands against multiple kots instances",
		Long:         `Run commands against multiple kots instances`,
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", filepath.Join(home, ".kotsd.yaml"),
		"config file (default $HOME/.kotsd.yaml)")

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(AddInstanceCmd())
	cmd.AddCommand(UpdateInstanceCmd())
	cmd.AddCommand(ListInstancesCmd())
	cmd.AddCommand(DeleteInstanceCmd())

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(UpdateCmd())
	cmd.AddCommand(RedeployCmd())
	cmd.AddCommand(RemoveCmd())

	viper.BindPFlags(cmd.Flags())
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}

func InitAndExecute(v string) {
	version = v
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	_, err := os.Stat(cfgFile)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Config not found, creating.")
		if _, err := os.Create(cfgFile); err != nil { // perm 0666
			log.Fatal(err)
		}
	}

	d, _ := kotsd.ReadConfig(cfgFile)
	runtime_conf, _ = kotsd.ParseConfig(d)
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
