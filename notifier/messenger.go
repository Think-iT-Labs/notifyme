package notifier

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

const endpoint = "https://notifyme.think-it.io/notify"

type Messenger struct {
	Token string `json:"token"`
	Notification
}

var ErrWrongToken = errors.New("Wrong token, please verify you have the right token in your ~/.notifyme file or ask the bot to get a new one")

func (m Messenger) Notify() error {
	var output bytes.Buffer
	json.NewEncoder(&output).Encode(m)
	log.Debugf("Messenger Payload:\n%s", output.Bytes())
	res, err := http.Post(endpoint, "application/json", bufio.NewReader(&output))
	if err != nil {
		return err
	}

	response, _ := httputil.DumpResponse(res, true)
	log.Debugf("Server response:\n%s", response)

	if res.StatusCode == http.StatusNotFound {
		return ErrWrongToken
	}
	return nil
}
