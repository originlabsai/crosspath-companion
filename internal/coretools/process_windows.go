//go:build windows

package coretools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// setSysProcAttr is a no-op on Windows (no process groups via Setpgid).
func setSysProcAttr(cmd *exec.Cmd) {
	// Windows doesn't support Setpgid
}

// terminateProcess kills a process on Windows (no SIGTERM support).
func terminateProcess(proc *os.Process) error {
	return proc.Kill()
}

// forceKillProcess kills a process on Windows.
func forceKillProcess(proc *os.Process) {
	proc.Kill()
}

// isProcessAlive checks if a process is still running on Windows.
// This is used in the stop flow right before force-killing, so a false
// positive is harmless (we'd just attempt an extra kill).
func isProcessAlive(proc *os.Process) bool {
	// On Windows, os.FindProcess always succeeds and Signal(0) is not
	// supported. Use tasklist to check if the PID exists.
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", proc.Pid), "/FO", "CSV", "/NH").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), fmt.Sprintf("%d", proc.Pid))
}

// killProcessGroupForCleanup kills a process tree on Windows using taskkill.
func killProcessGroupForCleanup(pid int) {
	exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid)).Run()
}

// isRunningAsRoot returns false on Windows (no UID concept).
func isRunningAsRoot() bool {
	return false
}

// killOrphanedChromeProcesses kills Chrome processes on Windows.
func killOrphanedChromeProcesses() {
	// Use taskkill to find and kill Chrome with the MCP profile
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq chrome.exe", "/FO", "CSV", "/NH")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	if strings.Contains(string(out), "chrome.exe") {
		exec.Command("taskkill", "/F", "/IM", "chrome.exe", "/FI", "WINDOWTITLE eq *chrome-devtools-mcp*").Run()
	}
}

// isChromeProcessRunning checks if any Chrome process is running on Windows.
func isChromeProcessRunning() bool {
	out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq chrome.exe", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "chrome.exe")
}
