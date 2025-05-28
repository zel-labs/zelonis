package validator

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	datadirDefaultKeyStore = "keystore"
)

type Config struct {
	PrivateKey string

	Name string `toml:"-"`

	UserIdent string `toml:",omitempty"`

	Version string `toml:"-"`

	DataDir string

	KeyStoreDir string `toml:",omitempty"`

	ExternalSigner string `toml:",omitempty"`

	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort   int
	GossipSeed []string
	GossipPort int
	Validator  bool
	Stake      float64
}

const (
	DefaultHTTPHost   = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort   = 8545        // Default TCP port for the HTTP RPC server
	DefaultGossipPort = 30331       // Default TCP port for the HTTP RPC server

)

var concensuSender = []byte{
	0x5a, 0x65, 0x6c, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31,
}

var DefaultSetup = Config{
	DataDir:  DefaultDataDir(),
	HTTPHost: DefaultHTTPHost,
	HTTPPort: DefaultHTTPPort,
	GossipSeed: []string{
		"584d4e5035536673386d6a4c475a74",
		"3773714c78475459583167654776",
	},
	GossipPort: DefaultGossipPort,
	Validator:  false,
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		switch runtime.GOOS {
		case "darwin":
			return filepath.Join(home, "Library", "Zelonis")
		case "windows":
			// We used to put everything in %HOME%\AppData\Roaming, but this caused
			// problems with non-typical setups. If this fallback location exists and
			// is non-empty, use it, otherwise DTRT and check %LOCALAPPDATA%.
			fallback := filepath.Join(home, "AppData", "Roaming", "Zelonis")
			appdata := windowsAppData()
			if appdata == "" || isNonEmptyDir(fallback) {
				return fallback
			}
			return filepath.Join(appdata, "Zelonis")
		default:
			return filepath.Join(home, ".Zelonis")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func windowsAppData() string {
	v := os.Getenv("LOCALAPPDATA")
	if v == "" {
		// Windows XP and below don't have LocalAppData. Crash here because
		// we don't support Windows XP and undefining the variable will cause
		// other issues.
		panic("environment variable LocalAppData is undefined")
	}
	return v
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func isNonEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	names, _ := f.Readdir(1)
	f.Close()
	return len(names) > 0
}
