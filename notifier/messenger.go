package notifier

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
)

const endpoint = "https://notifyme.think-it.io/notify"

type Messenger struct {
	Token string `json:"token"`
	Notification
}

func (m Messenger) Notify() {
	var output bytes.Buffer
	json.NewEncoder(&output).Encode(m)
	http.Post(endpoint, "application/json", bufio.NewReader(&output))
}
