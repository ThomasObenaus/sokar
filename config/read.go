package config

import "fmt"

func (cfg *Config) ReadConfig(args []string) error {
	if args == nil || len(args) == 0 {
		return fmt.Errorf("Args are missing")
	}

	if cfg.pFlagSet == nil {
		return fmt.Errorf("Pflag is nil")
	}

	if cfg.viper == nil {
		return fmt.Errorf("Viper is nil")
	}

	if err := cfg.setDefaults(); err != nil {
		return err
	}

	if err := cfg.registerAndParseFlags(args); err != nil {
		return err
	}

	cfgFile := cfg.viper.GetString(configFile.name)
	if err := cfg.readCfgFile(cfgFile); err != nil {
		return err
	}

	if err := cfg.registerEnvParams(); err != nil {
		return err
	}

	return cfg.fillCfgValues()
}

func (cfg *Config) readCfgFile(cfgFileName string) error {
	if len(cfgFileName) == 0 {
		return nil
	}
	cfg.viper.SetConfigFile(cfgFileName)
	if err := cfg.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
