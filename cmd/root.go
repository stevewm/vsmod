package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"vsmod/internal/api"
	"vsmod/internal/config"

	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var configFileLocation string
var conf config.ConfigFile
var modAPI = api.NewModAPI(&http.Client{Timeout: 10 * time.Second})

var rootCmd = &cobra.Command{
	Use:              "vsmod",
	Short:            "A CLI tool for managing Vintage Story mods",
	Long:             `vsmod is a CLI tool for managing Vintage Story mods in a declarative manner using a config file.`,
	PersistentPreRun: toggleDebug,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFileLocation, "file", "", "config file (default is $PWD/mods.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "verbose logging")
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func initConfig() {
	if configFileLocation != "" {
		viper.SetConfigFile(configFileLocation)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("mods")
	}

	d, err := os.ReadFile(configFileLocation)
	if err != nil {
		log.Errorf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.UnmarshalStrict(d, &conf); err != nil {
		log.Errorf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Errorf("Error unmarshalling config file: %v\n", err)
		os.Exit(1)
	}
}
