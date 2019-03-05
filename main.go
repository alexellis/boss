package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/namespaces"
	"github.com/crosbymichael/boss/api"
	v1 "github.com/crosbymichael/boss/api/v1"
	"github.com/crosbymichael/boss/cmd"
	"github.com/crosbymichael/boss/version"
	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "boss"
	app.Version = version.Version
	app.Usage = "run containers like a ross"
	app.Description = cmd.Banner
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in the logs",
		},
		cli.StringFlag{
			Name:   "agent",
			Usage:  "agent address",
			Value:  "0.0.0.0:1337",
			EnvVar: "BOSS_AGENT",
		},
		cli.StringFlag{
			Name:   "sentry-dsn",
			Usage:  "sentry DSN",
			EnvVar: "SENTRY_DSN",
		},
	}
	app.Before = func(clix *cli.Context) error {
		if clix.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if dsn := clix.GlobalString("sentry-dsn"); dsn != "" {
			raven.SetDSN(dsn)
			raven.DefaultClient.SetRelease(version.Version)
		}
		return nil
	}
	app.Commands = []cli.Command{
		checkpointCommand,
		createCommand,
		deleteCommand,
		getCommand,
		killCommand,
		listCommand,
		migrateCommand,
		pushCommand,
		restoreCommand,
		rollbackCommand,
		startCommand,
		stopCommand,
		updateCommand,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		raven.CaptureErrorAndWait(err, nil)
		os.Exit(1)
	}
}

func Context() context.Context {
	return namespaces.WithNamespace(context.Background(), v1.DefaultNamespace)
}

func Agent(clix *cli.Context) (*api.LocalAgent, error) {
	return api.Agent(clix.GlobalString("agent"))
}
