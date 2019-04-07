package config

import "strings"

func (cfg *Config) registerEnvParams() {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	cfg.viper.SetEnvKeyReplacer(replacer)

	for _, entry := range cfg.configEntries {
		registerEnv(cfg.viper, entry)
	}
}

func (cfg *Config) registerAndParseFlags(args []string) {

	for _, entry := range cfg.configEntries {
		registerFlag(cfg.pFlagSet, entry)
	}

	cfg.pFlagSet.Parse(args)
	cfg.viper.BindPFlags(cfg.pFlagSet)
}

func (cfg *Config) setDefaults() {
	for _, entry := range cfg.configEntries {
		setDefault(cfg.viper, entry)
	}
}
