package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

const (
	waitStep                = time.Second * 5
	defaultDownstreamBranch = "master"
)

// Plugin defines the deploy plugin parameters.
type Plugin struct {
	Server           string
	Token            string
	WaitRepo         string
	BuildNumber      int
	DownstreamRepo   string
	DownstreamBranch string
}

// Exec runs the plugin.
func (p *Plugin) Exec() error {
	if len(p.Server) == 0 {
		return errors.New("missing Drone server uri")
	}

	if len(p.Token) == 0 {
		return errors.New("missing Drone access token")
	}

	if len(p.WaitRepo) == 0 {
		return errors.New("missing wait repository")
	}

	if p.BuildNumber == 0 {
		return errors.New("missing build number")
	}

	if len(p.DownstreamRepo) == 0 {
		return errors.New("missing downstream repository")
	}

	if len(p.DownstreamBranch) == 0 {
		log.Printf("using default downstream branch: %q", defaultDownstreamBranch)
		p.DownstreamBranch = defaultDownstreamBranch
	}

	config := new(oauth2.Config)
	auth := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: p.Token,
		},
	)
	client := drone.NewClient(p.Server, auth)

	waitRepoOwner, waitRepoName, err := parseRepo(p.WaitRepo)
	if err != nil {
		return fmt.Errorf("invalid wait repository: %q", p.WaitRepo)
	}

	// get current job number
	curPid := currentJobPID()

	// wait for other build processes to be completed and success before deploy
	for {
		build, err := client.Build(waitRepoOwner, waitRepoName, p.BuildNumber)
		if err != nil {
			return err
		}

		var toDeploy = true
		for _, proc := range build.Procs {
			if proc.PID == curPid {
				continue
			}

			switch proc.State {
			case "success":
				continue
			case "failure":
				return fmt.Errorf("do not deploy, job failure: %s ", proc.PID)
			default:
				log.Printf("job %d is not completed, state: %s", proc.PID, proc.State)
				toDeploy = false
				break
			}
		}
		if toDeploy {
			break
		}
		time.Sleep(waitStep)
	}

	downstreamRepoOwner, downstreamRepoName, err := parseRepo(p.DownstreamRepo)
	if err != nil {
		return fmt.Errorf("invalid downstream repository: %q", p.WaitRepo)
	}

	lb, err := client.BuildLast(downstreamRepoOwner, downstreamRepoName, p.DownstreamBranch)
	if err != nil {
		return err
	}

	nb, err := client.BuildFork(downstreamRepoOwner, downstreamRepoName, lb.Number, nil)
	if err != nil {
		return err
	}

	log.Printf("starting build: %d for %s", nb.Number, p.DownstreamRepo)
	return nil
}

// get current job number. Return the job number if this plugin is
// currently running in a Drone build step, otherwrise return 0.
func currentJobPID() int {
	curPidStr := os.Getenv("DRONE_JOB_NUMBER")
	if len(curPidStr) == 0 {
		return 0
	}
	curPid, err := strconv.Atoi(curPidStr)
	if err != nil {
		log.Println(err)
		return 0
	}
	return curPid
}

// parseRepo returns user and repo from given string with format: owner/repo.
func parseRepo(str string) (user, repo string, err error) {
	var parts = strings.Split(str, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid or missing repository: %q", str)
		return
	}
	user = parts[0]
	repo = parts[1]
	return
}
