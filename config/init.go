package config

import (
	"fmt"
	"log"
	"strings"
)

const keyEnvPrefix = "SKR"
const keyNomadServerAddress = "nomad-server-address"
const keyCfgFile = "config-file"
const keyDryMode = "dry-run"

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

	cfg.setDefaults()

	cfg.registerAndParseFlags(args)

	cfgFile := cfg.viper.GetString(keyCfgFile)
	if err := cfg.readCfgFile(cfgFile); err != nil {
		return err
	}

	cfg.registerEnvParams()

	cfg.fillCfgValues()
	return nil
}

func (cfg *Config) fillCfgValues() {

	cfg.Port = cfg.viper.GetInt("port")
	cfg.DryRunMode = cfg.viper.GetBool(keyDryMode)

}

func (cfg *Config) registerEnvParams() {
	replacer := strings.NewReplacer("-", "_")
	cfg.viper.SetEnvKeyReplacer(replacer)

	cfg.viper.SetEnvPrefix(keyEnvPrefix)

	cfg.viper.BindEnv(keyNomadServerAddress)
	cfg.viper.BindEnv(keyDryMode)
	cfg.viper.BindEnv(keyCfgFile)
}

func (cfg *Config) registerAndParseFlags(args []string) {
	cfg.pFlagSet.String(keyNomadServerAddress, "", "Specifies the address of the nomad server.")
	cfg.pFlagSet.Bool(keyDryMode, false, "If true, then sokar won't execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed.")
	cfg.pFlagSet.String(keyCfgFile, "", "Specifies the full path and name of the configuration file for sokar.")

	cfg.pFlagSet.Parse(args)
	cfg.viper.BindPFlags(cfg.pFlagSet)
}

func (cfg *Config) setDefaults() {
	cfg.viper.SetDefault(keyDryMode, false)
	cfg.viper.SetDefault(keyNomadServerAddress, "http://localhost:4646")
	cfg.viper.SetDefault(keyCfgFile, "")
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
func InitMe(args []string) {

	cfg := NewDefaultConfig()
	err := cfg.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	log.Printf("NSRV %s", cfg.viper.Get(keyNomadServerAddress))

	////	pflag.NewFlagSet(name, errorHandling)
	//
	//pflag.BoolP("verbose", "v", false, "verbose output")
	//
	//viper.SetDefault("dry_run_mode", false)
	//
	//viper.SetConfigName("full")             // name of config file (without extension)
	//viper.AddConfigPath("examples/config/") // path to look for the config file in
	//err := viper.ReadInConfig()             // Find and read the config file
	//if err != nil {                         // Handle errors reading the config file
	//	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	//}
	//
	//log.Printf("DUR %s\n", viper.GetDuration("capacity_planner.down_scale_cooldown").String())
	//
	//viper.SetEnvPrefix("SKR") // will be uppercased automatically
	//viper.BindEnv("dry_run_mode")
	//
	////flag.Int("flagname", 1234, "help message for flagname")
	////
	////pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	////pflag.Parse()
	////viper.BindPFlags(pflag.CommandLine)
}
