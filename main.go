package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mackerelio/checkers"
	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
)

const (
	Name    string = "mkr-check"
	Version string = "0.1.1"
)

var statusToExitCode = map[checks.Status]checkers.Status{
	checks.StatusOK:       checkers.OK,
	checks.StatusWarning:  checkers.WARNING,
	checks.StatusCritical: checkers.CRITICAL,
	checks.StatusUnknown:  checkers.UNKNOWN,
}

func run(args []string) int {
	var (
		conffile string
		version  bool
	)
	flags := flag.NewFlagSet("mkr-check", flag.ContinueOnError)
	flags.StringVar(&conffile, "c", "/etc/mackerel-agent/mackerel-agent.conf", "mackerel-agent's config file")
	flags.StringVar(&conffile, "config", "/etc/mackerel-agent/mackerel-agent.conf", "mackerel-agent's config file")
	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	if err := flags.Parse(args); err != nil {
		log.Println(err)
		return 10
	}

	if version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		return 0
	}

	conf, err := config.LoadConfig(conffile)
	if err != nil {
		log.Println(err)
		return 11
	}

	agent := command.NewAgent(conf)
	sort.Slice(agent.Checkers, func(i, j int) bool {
		return strings.Compare(agent.Checkers[i].Name, agent.Checkers[j].Name) == -1
	})
	for _, c := range agent.Checkers {
		report := c.Check()
		var exitCode checkers.Status
		if code, ok := statusToExitCode[report.Status]; ok {
			exitCode = code
		}
		ckr := checkers.NewChecker(exitCode, report.Message)
		fmt.Printf("[%v]\n%v", c.Name, ckr.String())
		if report.Status != checks.StatusOK {
			return int(ckr.Status)
		}
	}

	return int(checkers.OK)
}

func main() {
	os.Exit(run(os.Args[1:]))
}
