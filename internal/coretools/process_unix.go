//go:build !windows

package coretools

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// setSysProcAttr sets process group ID so child processes can be killed as a group.
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup sends a signal to the entire process group.
func killProcessGroup(pid int, sig syscall.Signal) error {
	return syscall.Kill(-pid, sig)
}

// terminateProcess sends SIGTERM to a process.
func terminateProcess(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}

// forceKillProcess sends SIGKILL to a process.
func forceKillProcess(proc *os.Process) {
	proc.Signal(syscall.SIGKILL)
}

// isProcessAlive checks if a process is still running.
func isProcessAlive(proc *os.Process) bool {
	return proc.Signal(syscall.Signal(0)) == nil
}

// killProcessGroupForCleanup kills a process group with SIGTERM then SIGKILL after a timeout.
func killProcessGroupForCleanup(pid int) {
	syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(5 * time.Second)
	syscall.Kill(-pid, syscall.SIGKILL)
}

// isRunningAsRoot returns true if the process is running with UID 0.
func isRunningAsRoot() bool {
	return os.Getuid() == 0
}

// killOrphanedChromeProcesses kills Chrome processes matching MCP profile patterns.
func killOrphanedChromeProcesses() {
	patterns := []string{
		"chrome-devtools-mcp/chrome-profile",
	}
	for _, pattern := range patterns {
		cmd := exec.Command("pkill", "-f", pattern)
		cmd.Run()
	}
}

// isChromeProcessRunning checks if any Chrome process is running.
func isChromeProcessRunning() bool {
	out, _ := exec.Command("pgrep", "-x", "chrome").Output()
	return len(out) > 0
}
