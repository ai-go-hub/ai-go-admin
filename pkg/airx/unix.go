//go:build !windows

package airx

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// AirPID 返回父进程 PID（若父进程是 air）
func AirPID() (int, bool) {
	pid := os.Getppid()
	if pid <= 0 {
		return 0, false
	}
	if !isAirProcess(pid) {
		return 0, false
	}
	return pid, true
}

// isAirProcess 判断进程是否为 air（优先读 /proc，退回 ps 命令）
// MacOS 无 /proc，且 ps -o comm= 可能输出完整路径，统一取基名比较
func isAirProcess(pid int) bool {
	if name, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/comm"); err == nil {
		return filepath.Base(strings.TrimSpace(string(name))) == "air"
	}
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=").Output()
	if err != nil {
		return false
	}
	return filepath.Base(strings.TrimSpace(string(out))) == "air"
}

// pauseProcess 通过 SIGSTOP 挂起进程
func pauseProcess(pid int) {
	_ = syscall.Kill(pid, syscall.SIGSTOP)
}

// resumeProcess 通过 SIGCONT 恢复进程
func resumeProcess(pid int) {
	_ = syscall.Kill(pid, syscall.SIGCONT)
}
