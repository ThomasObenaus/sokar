package config

import (
	"strings"
)

func (cfg *Config) registerEnvParams() error {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	cfg.viper.SetEnvKeyReplacer(replacer)

	for _, entry := range cfg.configEntries {
		if err := registerEnv(cfg.viper, entry); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *Config) registerAndParseFlags(args []string) error {

	for _, entry := range cfg.configEntries {
		if err := registerFlag(cfg.pFlagSet, entry); err != nil {
			return err
		}
	}

	if err := cfg.pFlagSet.Parse(args); err != nil {
		return err
	}
	cfg.viper.BindPFlags(cfg.pFlagSet)

	return nil
}

func (cfg *Config) setDefaults() error {
	for _, entry := range cfg.configEntries {
		if err := setDefault(cfg.viper, entry); err != nil {
			return err
		}
	}
	return nil
}
