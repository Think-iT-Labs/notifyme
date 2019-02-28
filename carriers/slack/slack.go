package slack

import (
	"fmt"
	"strings"
	"sync"

	"github.com/imdario/mergo"
	"github.com/nlopes/slack"
	"github.com/think-it-labs/notifyme/carriers"
	"github.com/think-it-labs/notifyme/notification"
)

type Slack struct {
	Token    string
	Channels string
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
	// Default config
	var slackCarrierConfig = Slack{}

	if err := mergo.Map(&slackCarrierConfig, conf, mergo.WithOverride); err != nil {
		return nil, err
	}

	if slackCarrierConfig.Token == "" {
		return nil, fmt.Errorf("Slack: missing token")
	}

	// Build channel list
	for _, channel := range strings.Split(slackCarrierConfig.Channels, ",") {
		slackCarrierConfig.channels = append(slackCarrierConfig.channels, strings.TrimSpace(channel))
	}

	return &slackCarrierConfig, nil
}

// Send will send the notification to the desired channels.
func (c *Slack) Send(notif *notification.Notification) (err error) {
	api := slack.New(c.Token)

	postMessage := buildPostMessage(notif)

	// Send notifications
	var wg sync.WaitGroup
	wg.Add(len(c.channels))
	for _, channel := range c.channels {
		go func(channel string) {
			_, _, err = api.PostMessage(channel, postMessage)

			wg.Done()
		}(channel)

	}
	wg.Wait()

	return
}

func buildPostMessage(notif *notification.Notification) slack.MsgOption {

	body := fmt.Sprintf("```\n%s```", notif.Logs)
	cmd := fmt.Sprintf("$ %s", notif.Cmd)
	color := GREENCOLOR
	if notif.ExitCode != 0 {
		color = REDCOLOR
	}

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

	attachmentOpts := slack.MsgOptionAttachments(attachment)

	return slack.MsgOptionCompose(attachmentOpts)
}
