package carriers

import (
	"fmt"
	"sync"

	"github.com/nlopes/slack"
	"github.com/think-it-labs/notifyme/carriers"
	"github.com/think-it-labs/notifyme/notification"
)

type Slack struct {
	token    string
	channels []string
}

const (
	REDCOLOR   = "#ff0000"
	GREENCOLOR = "#009a8d"
)

func init() {
	carriers.RegisterCarrier("slack", new)
}

func new(conf map[string]interface{}) (carriers.Carrier, error) {
	token, ok := conf["token"]
	if !ok {
		return nil, fmt.Errorf("Slack: missing token")
	}

	// Build channels list
	var channels []string
	for _, channel := range conf["channels"].([]interface{}) {
		channels = append(channels, channel.(string))
	}
	return &Slack{
		token:    token.(string),
		channels: channels,
	}, nil
}

// Send will send the notification to the desired channels.
func (c *Slack) Send(notif *notification.Notification) error {
	// TODO
	api := slack.New(c.token)

	postMessage := buildPostMessage(notif)
	title := ""
	var err error

	// Send notifications
	var wg sync.WaitGroup
	wg.Add(len(c.channels))
	for _, channel := range c.channels {
		go func(channel string) {
			_, _, err = api.PostMessage(
				channel,
				title,
				postMessage,
			)
			wg.Done()
		}(channel)

	}
	wg.Wait()

	return err
}

func buildPostMessage(notif *notification.Notification) slack.PostMessageParameters {

	body := fmt.Sprintf("```\n%s```", notif.Logs)
	cmd := fmt.Sprintf("$ %s", notif.Cmd)
	color := GREENCOLOR
	if notif.ExitCode != 0 {
		color = REDCOLOR
	}
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{
		Color: color,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: cmd,
				Value: "",
			},
		},
		Text: body,
	}

	params.Attachments = []slack.Attachment{attachment}

	return params
}
