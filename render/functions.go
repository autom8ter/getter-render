package render

import (
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"text/template"
)

var functions = map[string]interface{}{
	"get":            viper.Get,
	"getString":      viper.GetString,
	"getStringSlice": viper.GetStringSlice,
	"getBool":        viper.GetBool,
	"getInt":         viper.GetInt,
	"allSettings":    viper.AllSettings,
	"allKeys":        viper.AllKeys,
}

func funcMap() template.FuncMap {
	for k, fn := range sprig.GenericFuncMap() {
		functions[k] = fn
	}
	return functions
}

func AddTmplFunctions(fns map[string]interface{}) {
	for k, v := range fns {
		functions[k] = v
	}
}
