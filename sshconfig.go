package herd

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type sshConfig []sshConfigBlock

type sshConfigBlock struct {
	glob   string
	config map[string]string
}

// Parse an openssh config file
func parseSshConfig(file string) (*sshConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Anything before the first Host section is global
	ret := sshConfig{}
	conf := make(map[string]string)
	hosts := []string{"*"}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		// Ignore comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			logrus.Errorf("Ignoring invalid ssh config line: %s", line)
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		// New host section
		if key == "host" {
			for _, h := range hosts {
				ret = append(ret, sshConfigBlock{glob: h, config: conf})
			}
			conf = make(map[string]string)
			hosts = strings.Split(val, " ")
			for i, h := range hosts {
				hosts[i] = strings.TrimSpace(h)
			}
			continue
		}
		// Can't overwrite things, first value rules
		if _, ok := conf[key]; !ok {
			conf[key] = val
		}
	}
	for _, h := range hosts {
		ret = append(ret, sshConfigBlock{glob: h, config: conf})
	}
	return &ret, nil
}

// Find all variables relevant for a host, first match wins
func (s *sshConfig) configForHost(name string) map[string]string {
	ret := make(map[string]string)
	for k, v := range puttyConfig(name) {
		ret[k] = v
	}
	if s == nil {
		return ret
	}
	for _, b := range *s {
		if ok, err := filepath.Match(b.glob, name); ok && err == nil {
			for k, v := range b.config {
				if _, ok := ret[k]; !ok {
					ret[k] = v
				}
			}
		}
	}
	for k, v := range puttyConfig(name) {
		ret[k] = v
	}
	return ret
}
