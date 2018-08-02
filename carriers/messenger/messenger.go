package messenger

import (
	"errors"
	"fmt"

	"github.com/think-it-labs/notifyme/carriers"
	"github.com/think-it-labs/notifyme/notification"
)

type Messenger struct {
	token string
}

func init() {
	carriers.RegisterCarrier("messenger", new)
}

func new(conf map[string]interface{}) (carriers.Carrier, error) {
	token, ok := conf["token"]
	if !ok {
		return nil, fmt.Errorf("Messenger: missing token")
	}
	return &Messenger{
		token: token.(string),
	}, nil
}

var errWrongToken = errors.New("Wrong token, please verify you have the right token in your ~/.notifyme file or ask the bot to get a new one")

// Send will send the notification using facebook messenger api
func (m *Messenger) Send(notif *notification.Notification) error {
	// TODO
	return nil
}
