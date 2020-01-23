// +build !darwin,!windows

package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func getCurrentUser() (*userData, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("I don't know who you are: %s", err)
	}
	if u.HomeDir == "" {
		return nil, fmt.Errorf("You don't have a homedir")
	}
	usr := userData{
		user:            u,
		cacheDir:        filepath.Join(u.HomeDir, ".cache", "katyusha"),
		configDir:       filepath.Join(u.HomeDir, ".config", "katyusha"),
		systemConfigDir: "/etc/katyusha",
		dataDir:         filepath.Join(u.HomeDir, ".local", "share", "katyusha"),
		historyDir:      filepath.Join(u.HomeDir, ".local", "share", "katyusha", "history"),
	}
	if d, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
		usr.dataDir = filepath.Join(d, "katyusha")
		usr.historyDir = filepath.Join(d, "katyusha", "history")
	}
	if d, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		usr.configDir = filepath.Join(d, "katyusha")
	}
	if d, ok := os.LookupEnv("XDG_CACHE_HOME"); ok {
		usr.cacheDir = filepath.Join(d, "katyusha")
	}
	return &usr, nil
}
