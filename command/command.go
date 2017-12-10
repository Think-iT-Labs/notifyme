package command

import (
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Command{
		Cmd: cmd,
	}
}

func (c *Command) AddStdoutWriter(writer io.Writer) {
	c.Stdout = io.MultiWriter(c.Stdout, writer)
}

func (c *Command) AddStderrWriter(writer io.Writer) {
	c.Stderr = io.MultiWriter(c.Stderr, writer)
}

func (c *Command) Start() error {
	err := c.Cmd.Start()
	return err
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
