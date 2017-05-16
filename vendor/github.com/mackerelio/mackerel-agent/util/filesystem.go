// +build linux darwin freebsd netbsd

package util

// ref. https://github.com/opscode/ohai/blob/master/lib/ohai/plugins/linux/filesystem.rb

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
	"github.com/mackerelio/mackerel-agent/logging"
)

// DfStat is disk free statistics from df command.
// Field names are taken from column names of `df -P`
type DfStat struct {
	Name      string
	Blocks    uint64
	Used      uint64
	Available uint64
	Capacity  uint8
	Mounted   string
}

// `df -P` sample:
//  Filesystem     1024-blocks     Used Available Capacity Mounted on
//  /dev/sda1           19734388 16868164 1863772  91% /
//  tmpfs                 517224        0  517224   0% /lib/init/rw
//  udev                  512780       96  512684   1% /dev
//  tmpfs                 517224        4  517220   1% /dev/shm

var dfHeaderPattern = regexp.MustCompile(
	// 1024-blocks or 1k-blocks
	`^Filesystem\s+(?:1024|1[Kk])-block`,
)

var logger = logging.GetLogger("util.filesystem")

var dfOpt []string

func init() {
	switch runtime.GOOS {
	case "darwin":
		dfOpt = []string{"-Pkl"}
	case "freebsd":
		dfOpt = []string{"-Pkt", "noprocfs,devfs,fdescfs,nfs,cd9660"}
	case "netbsd":
		dfOpt = []string{"-Pkl"}
	default:
		dfOpt = []string{"-P"}
	}
}

// CollectDfValues collects disk free statistics from df command
func CollectDfValues() ([]*DfStat, error) {
	cmd := exec.Command("df", dfOpt...)
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  15 * time.Second,
		KillAfter: 5 * time.Second,
	}
	// Ignores exit status in case that df returns exit status 1
	// when the agent does not have permission to access file system info.
	_, stdout, _, err := tio.Run()

	if err != nil {
		logger.Warningf("'df %s' command exited with a non-zero status: '%s'", strings.Join(dfOpt, " "), err)
		return nil, nil
	}
	return parseDfLines(stdout), nil
}

func parseDfLines(out string) []*DfStat {
	lineScanner := bufio.NewScanner(strings.NewReader(out))
	var filesystems []*DfStat
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if dfHeaderPattern.MatchString(line) {
			continue
		}
		dfstat, err := parseDfLine(line)
		if err != nil {
			logger.Warningf(err.Error())
			continue
		}
		// https://github.com/docker/docker/blob/v1.5.0/daemon/graphdriver/devmapper/deviceset.go#L981
		if strings.HasPrefix(dfstat.Name, "/dev/mapper/docker-") {
			continue
		}
		filesystems = append(filesystems, dfstat)
	}
	return filesystems
}

var dfColumnsPattern = regexp.MustCompile(`^(.+?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.+)$`)

func parseDfLine(line string) (*DfStat, error) {
	matches := dfColumnsPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse line: [%s]", line)
	}
	name := matches[1]
	blocks, _ := strconv.ParseUint(matches[2], 0, 64)
	used, _ := strconv.ParseUint(matches[3], 0, 64)
	available, _ := strconv.ParseUint(matches[4], 0, 64)
	capacity, _ := strconv.ParseUint(matches[5], 0, 8)
	mounted := matches[6]

	return &DfStat{
		Name:      name,
		Blocks:    blocks,
		Used:      used,
		Available: available,
		Capacity:  uint8(capacity),
		Mounted:   mounted,
	}, nil
}
