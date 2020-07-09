package render

import (
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"text/template"
)

func funcMap() template.FuncMap {
	functions := map[string]interface{}{
		"get":            viper.Get,
		"getString":      viper.GetString,
		"getStringSlice": viper.GetStringSlice,
		"getBool":        viper.GetBool,
		"getInt":         viper.GetInt,
		"allSettings":    viper.AllSettings,
		"allKeys":        viper.AllKeys,
	}
	for k, fn := range sprig.GenericFuncMap() {
		functions[k] = fn
	}
	return functions
}
