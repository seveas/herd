package main

import (
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registryRefreshCmd = &cobra.Command{
	Use:   "refresh [provider...]",
	Short: "Refresh the cache of the registry",
	Long:  `Refresh the cache of the registry. If no providers are specified, all caches will be refreshed.`,
	RunE:  runRegistryRefresh,
}

func init() {
	rootCmd.AddCommand(registryRefreshCmd)
}

func runRegistryRefresh(cmd *cobra.Command, args []string) error {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	viper.Set("Refresh", true)

	// Only refresh caches, and only for the specified providers, if any were specified
	providers := args
	allProviders := len(providers) == 0
	fullConf := viper.Sub("Providers")
	overrides := make(map[string]any)
	refreshed := make([]string, 0)
	for provider := range fullConf.AllSettings() {
		conf := fullConf.Sub(provider)
		if !conf.GetBool("enabled") {
			continue
		}
		if conf.GetString("Provider") != "cache" {
			overrides[provider] = map[string]any{"Enabled": false}
			continue
		}
		if !allProviders && !slices.Contains(providers, provider) {
			overrides[provider] = map[string]any{"Enabled": false}
			continue
		}
		refreshed = append(refreshed, provider)
	}
	if len(overrides) > 0 {
		overrides = map[string]any{"Providers": overrides}
		err := viper.MergeConfigMap(overrides)
		if err != nil {
			panic(err)
		}
	}

	engine, err := setupScriptEngine(nil)
	if err != nil {
		return err
	}
	defer engine.End()
	if len(refreshed) == 0 {
		logrus.Info("No caches to refresh")
	} else {
		logrus.Infof("Cache refreshed for %s", strings.Join(refreshed, ", "))
	}
	return nil
}
