package ssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/seveas/herd"
)

func cpTestdata(name string, destDir string) error {
	bf, err := ioutil.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		return err
	}
	dest := filepath.Join(destDir, name)
	err = ioutil.WriteFile(dest, bf, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func TestIncludes(t *testing.T) {
	testDir, err := os.MkdirTemp(os.TempDir(), "test-includes")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	configDir := filepath.Join(testDir, "config.d")
	err = os.Mkdir(configDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = cpTestdata("include_prod", configDir)
	if err != nil {
		t.Fatal(err)
	}
	err = cpTestdata("include_test", configDir)
	if err != nil {
		t.Fatal(err)
	}

	assertConfig := func(t *testing.T, bl []*configBlock) {
		if len(bl) != 9 {
			t.Errorf("expected %d blocks, but got %d", 9, len(bl))
		}
		c := &config{
			user: user.User{
				Username: "bob",
				HomeDir:  "/home/bob",
			},
			blocks: bl,
		}
		type tcase struct {
			host    string
			expUser string
			expPort int
			expId   string
		}
		for i, tc := range []*tcase{
			{host: "", expUser: "bob", expPort: 22, expId: ""},
			{host: "xyz", expUser: "bob", expPort: 22, expId: ""},
			{host: "ab8-001.prod.dom", expUser: "srv-prod", expPort: 2022, expId: "/home/bob/.ssh/keys.d/id_key_prod"},
			{host: "ab8-002.prod.dom", expUser: "srv-prod", expPort: 2022, expId: "/home/bob/.ssh/keys.d/id_key_prod"},
			{host: "cb8-003.prod.dom", expUser: "srv-prod", expPort: 2022, expId: "/home/bob/.ssh/keys.d/id_key_prod"},
			{host: "ab8-004.test.dom", expUser: "srv-test", expPort: 2023, expId: "/home/bob/.ssh/keys.d/id_key_test"},
			{host: "fr5-002.test.dom", expUser: "srv-test", expPort: 2023, expId: "/home/bob/.ssh/keys.d/id_key_test"},
		} {
			cb := c.forHost(&herd.Host{
				Name: tc.host,
			})
			if cb.clientConfig.User != tc.expUser {
				t.Errorf("case %d: expected user %s, but got %s", i, tc.expUser, cb.clientConfig.User)
			}
			if cb.port != tc.expPort {
				t.Errorf("case %d: expected port %d, but got %d", i, tc.expPort, cb.port)
			}
			if cb.identityFile != tc.expId {
				t.Errorf("case %d: expected id %s, but got %s", i, tc.expId, cb.identityFile)
			}
		}

	}

	t.Run("multiple", func(t *testing.T) {
		sconf := fmt.Sprintf("Include %s %s",
			filepath.Join("config.d", "include_prod"),
			filepath.Join("config.d", "include_test"))
		conf := filepath.Join(testDir, "include_multiple")
		err = ioutil.WriteFile(conf, []byte(sconf), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		bl, err := parseConfig(true, testDir, "include_multiple")
		if err != nil {
			t.Fatal(err)
		}
		assertConfig(t, bl)
	})

	t.Run("simple", func(t *testing.T) {
		sconf := "Include config.d/*"
		conf := filepath.Join(testDir, "config")
		err = ioutil.WriteFile(conf, []byte(sconf), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		bl, err := parseConfig(true, testDir, "config")
		if err != nil {
			t.Fatal(err)
		}
		assertConfig(t, bl)
	})
}
