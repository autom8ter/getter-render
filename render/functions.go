package render

import (
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

func init() {
	for k, fn := range sprig.GenericFuncMap() {
		functions[k] = fn
	}
}

var functions = map[string]interface{}{
	"get":                viper.Get,
	"getString":          viper.GetString,
	"getStringSlice":     viper.GetStringSlice,
	"getBool":            viper.GetBool,
	"getInt":             viper.GetInt,
	"getStringMap":       viper.GetStringMap,
	"getStringMapString": viper.GetStringMapString,
	"getIntSlice":        viper.GetIntSlice,
	"getTime":            viper.GetTime,
	"getFloat64":         viper.GetFloat64,
	"getDuration":        viper.GetDuration,
	"isSet":              viper.IsSet,
	"inConfig":           viper.InConfig,
	"allKeys":            viper.AllKeys,
	"allSettings":        viper.AllSettings,
	"getEnv":             os.Getenv,
}

func funcMap() template.FuncMap {
	return functions
}

func AddTmplFunctions(fns map[string]interface{}) {
	for k, v := range fns {
		functions[k] = v
	}
}
