package cmd

import (
	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/scripting"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupScriptEngine() (*scripting.ScriptEngine, error) {
	formatter := katyusha.Formatters[viper.GetString("Formatter")]
	ui := katyusha.NewSimpleUI(formatter)
	ui.SetOutputMode(viper.Get("Output").(katyusha.OutputMode))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	logrus.SetFormatter(formatter)
	logrus.SetOutput(ui)

	registry, err := katyusha.NewRegistry()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	err = registry.Load(ui.CacheUpdateChannel())
	if err != nil {
		// Do not log this error, registry.Load() does its own error logging
		return nil, err
	}
	return scripting.NewScriptEngine(ui, katyusha.NewRunner(registry)), nil
}
