package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"
)

const (
	defaultExitCode = 1
)

type Command struct {
	*exec.Cmd
	waitErr error
}

func New(args []string) *Command {
	cmd := exec.Command(args[0], args[1:]...)

	return &Command{
		Cmd: cmd,
	}
}

func (c *Command) Start() (*bytes.Buffer, error) {
	output := new(bytes.Buffer)
	c.Cmd.Stdin = os.Stdin
	c.Cmd.Stdout = io.MultiWriter(os.Stdout, output)
	c.Cmd.Stderr = io.MultiWriter(os.Stderr, output)
	err := c.Cmd.Start()
	return output, err
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
