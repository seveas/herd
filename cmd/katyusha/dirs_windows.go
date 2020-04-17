package main

import (
	"fmt"
	"os/user"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func getCurrentUser() (*userData, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("I don't know who you are: %s", err)
	}
	usr := &userData{
		user: u,
	}

	d, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, windows.KF_FLAG_CREATE)
	if err != nil {
		return nil, err
	}
	usr.cacheDir = filepath.Join(d, "katyusha", "cache")
	usr.historyDir = filepath.Join(d, "katyusha", "history")

	d, err = windows.KnownFolderPath(windows.FOLDERID_RoamingAppData, windows.KF_FLAG_CREATE)
	if err != nil {
		return nil, err
	}
	usr.configDir = filepath.Join(d, "katyusha")
	usr.dataDir = filepath.Join(d, "katyusha")

	d, err = windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_CREATE)
	if err != nil {
		return nil, err
	}
	usr.systemConfigDir = filepath.Join(d, "katyusha")

	return usr, nil
}
