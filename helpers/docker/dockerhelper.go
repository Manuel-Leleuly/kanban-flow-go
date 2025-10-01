package dockerhelper

import (
	"os"
	"strings"
)

func IsInDocker() bool {
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "/docker/") || strings.Contains(string(data), "/lxc/")
}
