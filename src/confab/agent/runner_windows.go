package agent

import (
	"os"
	"syscall"
)

func signalProcess(proc *os.Process) error {
	err := proc.Signal(syscall.Signal(0))
	if err == syscall.EWINDOWS {
		return nil
	}
	return err
}
