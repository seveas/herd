package cmd

import (
	"github.com/seveas/herd"
	"github.com/seveas/herd/scripting"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupScriptEngine() (*scripting.ScriptEngine, error) {
	formatter := herd.Formatters[viper.GetString("Formatter")]
	ui := herd.NewSimpleUI(formatter)
	ui.SetOutputMode(viper.Get("Output").(herd.OutputMode))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	logrus.SetFormatter(formatter)
	logrus.SetOutput(ui)

	registry, err := herd.NewRegistry()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	err = registry.Load(ui.CacheUpdateChannel())
	if err != nil {
		// Do not log this error, registry.Load() does its own error logging
		return nil, err
	}
	runner := herd.NewRunner(registry)
	runner.SetParallel(viper.GetInt("Parallel"))
	runner.SetTimeout(viper.GetDuration("Timeout"))
	runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	return scripting.NewScriptEngine(ui, runner), nil
}
