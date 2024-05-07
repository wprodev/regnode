package cli

import (
	"strconv"

	"github.com/regnode/internal/node"
	"github.com/regnode/internal/rancher"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().Bool("worker", false, "Node will be registered as a worker node")
	registerCmd.Flags().Bool("etcd", false, "Node will be registered as a etcd node")
	registerCmd.Flags().Bool("controlplane", false, "Node will be registered as a controlplane node")
	registerCmd.Flags().IntP("continuous", "n", 3, "Rerun registration process if it fails, waits for <value> seconds after each try. Retry will only run if config is valid, if parameters ale malformed (e.g. URL) it fails with exit code 1.")

	viper.BindPFlag("worker", registerCmd.Flags().Lookup("worker"))
	viper.BindPFlag("etcd", registerCmd.Flags().Lookup("etcd"))
	viper.BindPFlag("controlplane", registerCmd.Flags().Lookup("controlplane"))
	viper.BindPFlag("continuous", registerCmd.Flags().Lookup("continuous"))

}

var registerCmd = &cobra.Command{
	Use:     "register",
	Example: "regnode register --etcd --controlplane cluster-1",
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Register host in the rancher cluster.",
	Long: `Registers node in the rancher cluster by calling rancher's HTTP API. 
Retrieves the node command and executes command in shell. 
Node command installs rancher agent and runs it with defined flags based on --roles flag.
`,
	Run:    runRegisterCmd,
	PreRun: resolveArgsAndFlags,
}

func runRegisterCmd(cmd *cobra.Command, args []string) {

	if !cfg.ClusterConfig.IsControlplane && !cfg.ClusterConfig.IsEtcd && !cfg.ClusterConfig.IsWorker {
		log.Fatal().Msg("Node not configured as any of: wroker, etcd, controlplane. Provide values in config or as an argument.")
	}
	log.Info().Msgf("Registering node in rancher cluster: %s", cfg.ClusterConfig.Name)

	rancherHttpClient := rancher.HttpClient(cfg.ApiURL, cfg.ApiToken, true)
	nodeCommand, err := rancherHttpClient.RetrieveNodeCommand(cfg.ClusterConfig.Name)
	if err != nil {
		log.Fatal().Msgf("Failed when retrieving node command: %s", err.Error())
	}
	log.Info().Msgf("Retrieved register node command from rancher: %s", nodeCommand)
	var shellOutput string
	shellOutput, err = node.RegisterWithCommand(nodeCommand, cfg.ClusterConfig.IsControlplane, cfg.ClusterConfig.IsEtcd, cfg.ClusterConfig.IsWorker)
	if err != nil {
		log.Fatal().Msgf("Error while executing shell command: %s, error: %s", shellOutput, err)
	}
	log.Debug().Msgf("Shell cmd output: %s", shellOutput)
}

// resolveArgsAndFlags merges flags with exiting config to enforce loading order for all parameters
func resolveArgsAndFlags(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		cfg.ClusterConfig.Name = args[0]
	} else if cfg.ClusterConfig.Name == "" {
		log.Fatal().Msg("Cluster name not found. Provide cluster name in config or as an argument.")
	}

	if flag := cmd.Flags().Lookup("worker"); flag.Changed {
		cfg.ClusterConfig.IsWorker = true
	}

	if flag := cmd.Flags().Lookup("etcd"); flag.Changed {
		cfg.ClusterConfig.IsEtcd = true
	}

	if flag := cmd.Flags().Lookup("controlplane"); flag.Changed {
		cfg.ClusterConfig.IsControlplane = true
	}

	if flag := cmd.Flags().Lookup("continuous"); flag.Changed {
		i, ok := strconv.Atoi(flag.Value.String())
		if ok != nil {
			log.Fatal().Msgf("Failed to parse flag --continuous with value %s", flag.Value.String())
		}
		cfg.Continuous = i
	}

	log.Debug().Msgf("Config merged: %v", cfg)
}
