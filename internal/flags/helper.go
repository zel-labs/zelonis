package flags

import (
	"github.com/urfave/cli"
	"zelonis/params"
	"zelonis/utils/version"
)

// NewApp creates an app with sane defaults.
func NewApp(usage string) *cli.App {
	git, _ := version.VCS()
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = params.VersionWithCommit(git.Commit, git.Date)
	app.Usage = usage
	app.Copyright = "Copyright 2025-2026 The go-Zelonis Authors"

	return app
}
