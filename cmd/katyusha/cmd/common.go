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

	registry := katyusha.NewRegistry(viper.GetString("RootDir"))
	registry.SetSortFields(viper.GetStringSlice("Sort"))
	registry.LoadMagicProviders()
	conf := viper.Sub("Providers")
	if conf != nil {
		if err := registry.LoadProviders(conf); err != nil {
			logrus.Error(err.Error())
			ui.End()
			return nil, err
		}
	}
	if err := registry.LoadHosts(ui.CacheUpdateChannel()); err != nil {
		// Do not log this error, registry.LoadHosts() does its own error logging
		ui.End()
		return nil, err
	}
	runner := katyusha.NewRunner(registry)
	runner.SetParallel(viper.GetInt("Parallel"))
	runner.SetTimeout(viper.GetDuration("Timeout"))
	runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	return scripting.NewScriptEngine(ui, runner), nil
}
