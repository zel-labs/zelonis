package validator

import (
	"github.com/gofrs/flock"
	"os"
	"path/filepath"
	"strings"
	"zelonis/internal/flags"
)

var (
	KeyStoreDirFlag = &flags.DirectoryFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}
)

func (vn *Validator) openDataDir() error {
	if vn.cfg.DataDir == "" {
		return nil // ephemeral
	}

	instdir := filepath.Join(vn.cfg.DataDir, vn.cfg.name())
	if err := os.MkdirAll(instdir, 0700); err != nil {
		return err
	}
	vn.dirLock = flock.New(filepath.Join(instdir, "LOCK"))

	if locked, err := vn.dirLock.TryLock(); err != nil {
		return err
	} else if !locked {
		return ErrDatadirUsed
	}
	return nil

}

func (c *Config) name() string {
	if c.Name == "" {
		progname := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if progname == "" {
			panic("empty executable name, set Config.Name")
		}
		return progname
	}
	return c.Name
}
