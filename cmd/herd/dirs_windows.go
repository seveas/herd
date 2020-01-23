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
	usr.cacheDir = filepath.Join(d, "herd", "cache")
	usr.historyDir = filepath.Join(d, "herd", "history")

	d, err = windows.KnownFolderPath(windows.FOLDERID_RoamingAppData, windows.KF_FLAG_CREATE)
	if err != nil {
		return nil, err
	}
	usr.configDir = filepath.Join(d, "herd")
	usr.dataDir = filepath.Join(d, "herd")

	d, err = windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_CREATE)
	if err != nil {
		return nil, err
	}
	usr.systemConfigDir = filepath.Join(d, "herd")

	return usr, nil
}
