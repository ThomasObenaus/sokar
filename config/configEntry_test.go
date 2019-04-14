package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CheckViper(t *testing.T) {

	err := checkViper(nil, configEntry{})
	assert.Error(t, err)

	vp := viper.New()
	require.NotNil(t, vp)

	err = checkViper(vp, configEntry{})
	assert.Error(t, err)

	cfgE := configEntry{
		name: "bla",
	}
	err = checkViper(vp, cfgE)
	assert.NoError(t, err)
}

func Test_SetDefault_OK(t *testing.T) {

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := configEntry{
		name:         "bla",
		defaultValue: 20,
	}
	err := setDefault(vp, cfgE)
	assert.NoError(t, err)

	assert.NotNil(t, vp.GetInt(cfgE.name))
	assert.Equal(t, cfgE.defaultValue, vp.GetInt(cfgE.name))
}
func Test_SetDefault_Fail(t *testing.T) {

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := configEntry{
		name: "bla",
	}
	err := setDefault(vp, cfgE)
	assert.Error(t, err)
}
func Test_RegisterEnv_OK(t *testing.T) {

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := configEntry{
		name:    "flag",
		bindEnv: true,
	}
	err := registerEnv(vp, cfgE)
	assert.NoError(t, err)
	os.Setenv(EnvPrefix+"_"+strings.ToUpper(cfgE.name), "test")
	assert.NotEmpty(t, vp.Get(cfgE.name))

	cfgE = configEntry{
		name:    "flag",
		bindEnv: true,
	}
	err = registerEnv(vp, cfgE)
	assert.NoError(t, err)
	os.Setenv(strings.ToUpper(EnvPrefix+"_"+cfgE.name), "test")
	assert.NotEmpty(t, vp.Get(cfgE.name))
}

func Test_RegisterEnv_Fail(t *testing.T) {

	err := registerEnv(nil, configEntry{})
	assert.NoError(t, err)

	err = registerEnv(nil, configEntry{bindEnv: true})
	assert.Error(t, err)

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := configEntry{bindEnv: true}
	err = registerEnv(vp, cfgE)
	assert.Error(t, err)
}

func Test_RegisterFlag_Fail(t *testing.T) {

	err := registerFlag(nil, configEntry{})
	assert.NoError(t, err)

	err = registerFlag(nil, configEntry{bindFlag: true})
	assert.Error(t, err)

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	require.NotNil(t, flagSet)

	cfgE := configEntry{bindFlag: true}
	err = registerFlag(flagSet, cfgE)
	assert.Error(t, err)

	cfgE.name = "flag1"
	err = registerFlag(flagSet, cfgE)
	assert.Error(t, err)
}

func Test_RegisterFlag_Ok(t *testing.T) {

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	require.NotNil(t, flagSet)

	// String
	cfgE := configEntry{
		bindFlag:      true,
		name:          "flag1",
		defaultValue:  "default",
		usage:         "The default value",
		flagShortName: "a",
	}
	err := registerFlag(flagSet, cfgE)
	assert.NoError(t, err)
	flag := flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, cfgE.defaultValue.(string), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Uint
	cfgE = configEntry{
		bindFlag:      true,
		name:          "flag2",
		defaultValue:  uint(1),
		usage:         "An uint",
		flagShortName: "b",
	}
	err = registerFlag(flagSet, cfgE)
	assert.NoError(t, err)
	flag = flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, fmt.Sprintf("%d", cfgE.defaultValue.(uint)), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Int
	cfgE = configEntry{
		bindFlag:      true,
		name:          "flag3",
		defaultValue:  int(1),
		usage:         "An int",
		flagShortName: "c",
	}
	err = registerFlag(flagSet, cfgE)
	assert.NoError(t, err)
	flag = flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, fmt.Sprintf("%d", cfgE.defaultValue.(int)), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Type not supported
	typeNotSupported := reflect.TypeOf("")

	cfgE = configEntry{
		bindFlag:      true,
		name:          "flag4",
		defaultValue:  typeNotSupported,
		usage:         "Reflect type info",
		flagShortName: "d",
	}
	err = registerFlag(flagSet, cfgE)
	assert.Error(t, err)
}
