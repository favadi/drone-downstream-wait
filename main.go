package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var build = "0" // build number set at compile time

func main() {
	app := cli.NewApp()
	app.Name = "wait downstream plugin"
	app.Usage = "wait until all other jobs finished with success state " +
		"and trigger downstream repository build"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "Trigger a drone build on a custom server",
			EnvVar: "DEPLOY_SERVER,PLUGIN_SERVER,CI_SYSTEM_LINK",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Drone API token from your user setting",
			EnvVar: "DEPLOY_TOKEN,PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "wait_repository",
			Usage:  "Repository to wait before trigger downstream",
			EnvVar: "DEPLOY_WAIT_REPOSITORY,PLUGIN_WAIT_REPOSITORY,DRONE_REPO",
		},
		cli.IntFlag{
			Name:   "build",
			Usage:  "Drone build number to wait before trigger downstream",
			EnvVar: "DEPLOY_BUILD_NUMBER,PLUGIN_BUILD_NUMBER,DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "downstream_repository",
			Usage:  "Repository to trigger",
			EnvVar: "DEPLOYMENT_DOWNSTREAM_REPOSITORY,PLUGIN_DOWNSTREAM_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "downstream_branch",
			Usage:  "Branch of downstream repository",
			EnvVar: "DEPLOYMENT_DOWNSTREAM_BRANCH,PLUGIN_DOWNSTREAM_BRANCH",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	p := Plugin{
		Server:           c.String("server"),
		Token:            c.String("token"),
		WaitRepo:         c.String("wait_repository"),
		BuildNumber:      c.Int("build"),
		DownstreamRepo:   c.String("downstream_repository"),
		DownstreamBranch: c.String("downstream_branch"),
	}
	return p.Exec()
}
