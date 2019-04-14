package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type configEntry struct {
	name         string
	usage        string
	defaultValue interface{}

	bindFlag      bool
	flagShortName string

	bindEnv bool
}

// EnvPrefix is the prefix used for sokars environment variables
const EnvPrefix = "SK"

func checkViper(vp *viper.Viper, cfgEntry configEntry) error {
	if vp == nil {
		return fmt.Errorf("Viper is nil")
	}

	if len(cfgEntry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	return nil
}

func registerFlag(flagSet *pflag.FlagSet, cfgEntry configEntry) error {
	if !cfgEntry.bindFlag {
		return nil
	}
	if flagSet == nil {
		return fmt.Errorf("FlagSet is nil")
	}
	if len(cfgEntry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	if cfgEntry.defaultValue == nil {
		return fmt.Errorf("Default Value is missing")
	}

	switch cfgEntry.defaultValue.(type) {
	case string:
		flagSet.StringP(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(string), cfgEntry.usage)
	case uint:
		flagSet.UintP(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(uint), cfgEntry.usage)
	case int:
		flagSet.IntP(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(int), cfgEntry.usage)
	case bool:
		flagSet.BoolP(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(bool), cfgEntry.usage)
	case time.Duration:
		flagSet.DurationP(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(time.Duration), cfgEntry.usage)
	case float64:
		flagSet.Float64P(cfgEntry.name, cfgEntry.flagShortName, cfgEntry.defaultValue.(float64), cfgEntry.usage)
	default:
		return fmt.Errorf("Type %s is not yet supported", reflect.TypeOf(cfgEntry.defaultValue))
	}

	return nil
}

func setDefault(vp *viper.Viper, cfgEntry configEntry) error {
	if err := checkViper(vp, cfgEntry); err != nil {
		return err
	}
	if cfgEntry.defaultValue == nil {
		return fmt.Errorf("Default Value is missing")
	}
	vp.SetDefault(cfgEntry.name, cfgEntry.defaultValue)

	return nil
}

func registerEnv(vp *viper.Viper, cfgEntry configEntry) error {
	if !cfgEntry.bindEnv {
		return nil
	}
	if err := checkViper(vp, cfgEntry); err != nil {
		return err
	}

	vp.SetEnvPrefix(EnvPrefix)
	return vp.BindEnv(cfgEntry.name)
}
