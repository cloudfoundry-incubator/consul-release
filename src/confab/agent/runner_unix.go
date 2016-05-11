// +build !windows

package agent

import (
	"os"
	"syscall"
)

func signalProcess(proc *os.Process) error {
	return proc.Signal(syscall.Signal(0))
}
