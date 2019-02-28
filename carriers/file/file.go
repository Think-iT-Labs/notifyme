package file

import (
	"fmt"

	"github.com/think-it-labs/notifyme/notification"

	"github.com/imdario/mergo"
	"github.com/think-it-labs/notifyme/carriers"
)

// FileCarrier carry the notification to a file, named FileCarrier to avoid any confusion with the File type.
type FileCarrier struct {
	Filename  string
	Overwrite bool
	Append    bool
	Base64    bool
}

func init() {
	carriers.RegisterCarrier("file", new)
}

func new(conf map[string]interface{}) (carriers.Carrier, error) {
	// Default config
	var fileCarrierConfig = FileCarrier{
		Filename:  "{{ argv0 }}",
		Overwrite: true,
		Append:    false,
		Base64:    false,
	}

	if err := mergo.Map(&fileCarrierConfig, conf, mergo.WithOverride); err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", fileCarrierConfig)
	return &fileCarrierConfig, nil
}

func (c *FileCarrier) Send(notif *notification.Notification) error {
	return nil
}
