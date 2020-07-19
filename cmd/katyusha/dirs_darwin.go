package main

import (
	"fmt"
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
	return &userData{
		user:            u,
		cacheDir:        filepath.Join(u.HomeDir, "Library", "Caches", "katyusha"),
		configDir:       filepath.Join(u.HomeDir, "Library", "Preferences", "katyusha"),
		systemConfigDir: "/etc/katyusha",
		dataDir:         filepath.Join(u.HomeDir, "Library", "ApplicationSupport", "katyusha"),
		historyDir:      filepath.Join(u.HomeDir, "Library", "ApplicationSupport", "katyusha", "history"),
	}, nil
}
