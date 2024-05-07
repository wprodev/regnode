package cli

import (
	"os"
	"strings"

	"github.com/regnode/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileName string = ".regnode"
	configFileType string = "yaml"
	envVarsPrefix  string = "REGNODE"
)

var (
	cfgFilePath string
	cfg         config.Config

	rootCmd = &cobra.Command{
		Use:              "regnode",
		Short:            "Tool to register node in a Rancher cluster",
		Long:             `Regnode registers host in specified rancher cluster. It retrieves a node command from Rancher API and runs it on the host. The node command installs rancher agent and executes it with necessary flags so the node can join as a worker, etcd or controlplane node.`,
		Run:              runRootCmd,
		PersistentPreRun: resolveConfig,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "", "Config file YAML")
	rootCmd.PersistentFlags().StringP("api-url", "u", "", "Rancher API URL")
	rootCmd.PersistentFlags().StringP("api-token", "t", "", "Rancher API Token")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", false, "Rancher API Token")

	viper.BindPFlag("api-url", rootCmd.Flags().Lookup("api-url"))
	viper.BindPFlag("api-token", rootCmd.Flags().Lookup("api-token"))
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))

	rootCmd.MarkFlagsRequiredTogether("api-url", "api-token")

	//cobra.OnInitialize(initConfig)
}

// Execute the main logic of the command
func runRootCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// Runs before every executed command to handle variables load order
func resolveConfig(cmd *cobra.Command, args []string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if flag := cmd.Flags().Lookup("debug"); flag.Changed {
		if flag.Value.String() == "true" || flag.Value.String() == "1" {
			cfg.Debug = true
		}
	}

	if cfgFilePath != "" {
		viper.SetConfigFile(cfgFilePath)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(configFileName)
	}
	viper.SetConfigType(configFileType)
	viper.SetEnvPrefix(envVarsPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal().Msgf("Config found but could not load: %s", err)
		}
		log.Info().Msgf("Config not used.")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Msgf("Unable to decode %s into config struct, %v", viper.ConfigFileUsed(), err)
	}

	log.Info().Msgf("Using config: %s", viper.ConfigFileUsed())

	if flag := cmd.Flags().Lookup("api-url"); flag.Changed {
		cfg.ApiURL = flag.Value.String()
	}

	if flag := cmd.Flags().Lookup("api-token"); flag.Changed {
		cfg.ApiToken = flag.Value.String()
	}

	if cfg.Debug {
		log.Debug().Msgf("Values read from config: %v", cfg)
	}
}
