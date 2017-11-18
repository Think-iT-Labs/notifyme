package notifier

type Notification struct {
	Cmd      string `json:"cmd"`
	ExitCode int    `json:"exit_code"`
	Logs     []byte `json:"logs"`
}

type Notifier interface {
	Notify() error
}
