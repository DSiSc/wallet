package main

import (
	"fmt"
	"github.com/DSiSc/wallet/cmd"
	"github.com/DSiSc/wallet/utils"
	"github.com/urfave/cli"
	"os"
	"sort"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the wallet command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		utils.PasswordFileFlag,
		utils.LightKDFFlag,
	}

	rpcFlags     = []cli.Flag{}
	whisperFlags = []cli.Flag{}
	metricsFlags = []cli.Flag{}
)

func init() {
	app.Action = wallet
	app.HideVersion = true
	app.Copyright = "Copyright 2018-2023 The justitia Authors"
	app.Commands = []cli.Command{
		cmd.AccountCommand,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)

	app.Before = func(ctx *cli.Context) error {
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		//debug.Exit()
		//console.Stdin.Close()
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func wallet(ctx *cli.Context) error {
	fmt.Print("***wallet()")
	return nil
}
