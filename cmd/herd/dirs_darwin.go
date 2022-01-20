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
	if d, ok := os.LookupEnv("HOME"); ok {
		u.HomeDir = d
	}
	if u.HomeDir == "" {
		return nil, fmt.Errorf("You don't have a homedir")
	}
	return &userData{
		user:            u,
		cacheDir:        filepath.Join(u.HomeDir, "Library", "Caches", "herd"),
		configDir:       filepath.Join(u.HomeDir, "Library", "Preferences", "herd"),
		systemConfigDir: "/etc/herd",
		dataDir:         filepath.Join(u.HomeDir, "Library", "Application Support", "herd"),
		historyDir:      filepath.Join(u.HomeDir, "Library", "Application Support", "herd", "history"),
	}, nil
}
