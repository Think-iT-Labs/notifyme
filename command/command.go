package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/kr/pty"
)

const (
	defaultExitCode = 1
)

type Command struct {
	*exec.Cmd
	waitErr error
}

func New(args []string) Command {
	cmd := exec.Command(args[0], args[1:]...)

	return Command{
		Cmd: cmd,
	}
}

func (c *Command) StartWithTTY() (*os.File, error) {
	return pty.Start(c.Cmd)
}

func (c *Command) Wait() int {
	c.waitErr = c.Cmd.Wait()
	return c.exitCode()
}

func (c *Command) exitCode() int {
	err := c.waitErr
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			return ws.ExitStatus()
		}
		return defaultExitCode
	}
	ws := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	return ws.ExitStatus()
}
