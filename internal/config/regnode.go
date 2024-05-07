package config

// Configuration struct for storing environment variables and flag overrides

type Config struct {
	Debug         bool          `mapstructure:"debug"`
	ApiURL        string        `mapstructure:"api-url"`
	ApiToken      string        `mapstructure:"api-token"`
	ClusterConfig ClusterConfig `mapstructure:"cluster"`
	Continuous    int           `mapstructure:"continuous"`
}

type ClusterConfig struct {
	Name           string `mapstructure:"name"`
	IsWorker       bool   `mapstructure:"worker"`
	IsControlplane bool   `mapstructure:"controlplane"`
	IsEtcd         bool   `mapstructure:"etcd"`
}
