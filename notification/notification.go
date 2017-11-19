package notification

type NotificationData struct {
	Cmd      string `json:"cmd"`
	ExitCode int    `json:"exit_code"`
	Logs     []byte `json:"logs"`
}

type Notification interface {
	Send() error
}
