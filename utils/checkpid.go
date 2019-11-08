package utils

import (
	"fmt"
	"os"
	"syscall"

	"github.com/facebookgo/pidfile"
)

func CheckPid(pidfileName string) error {
	pidfile.SetPidfilePath(pidfileName)
	pid, err := pidfile.Read()
	if err == nil && pid > 0 {
		process, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("failed to find process: %v", err)
		}
		err = process.Signal(syscall.Signal(0))
		if err == nil {
			return fmt.Errorf("perhaps another lachesis is already running with pid %d", pid)
		}
	}

	if err := pidfile.Write(); err != nil {
		return fmt.Errorf("error writing into pidfile: %v", err)
	}

	return nil
}
