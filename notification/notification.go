package notification

type Notification struct {
	Cmd      string `json:"cmd"`
	ExitCode int    `json:"exit_code"`
	Logs     []byte `json:"logs"`
}

func New(cmd string, exitCode int, logs []byte) *Notification {
	return &Notification{
		Cmd:      cmd,
		ExitCode: exitCode,
		Logs:     logs,
	}
}
