package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	pflag.BoolP("verbose", "v", false, "verbose output")

	viper.SetDefault("dry_run_mode", false)

	viper.SetConfigName("full")             // name of config file (without extension)
	viper.AddConfigPath("examples/config/") // path to look for the config file in
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Printf("DUR %s\n", viper.GetDuration("capacity_planner.down_scale_cooldown").String())

	viper.SetEnvPrefix("SKR") // will be uppercased automatically
	viper.BindEnv("dry_run_mode")

	//flag.Int("flagname", 1234, "help message for flagname")
	//
	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//pflag.Parse()
	//viper.BindPFlags(pflag.CommandLine)
}
