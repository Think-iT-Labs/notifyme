package carriers

import (
	"fmt"

	"github.com/think-it-labs/notifyme/notification"
)

type Carrier interface {
	Send(*notification.Notification) error
}

type carrierInitFunc func(conf map[string]interface{}) (Carrier, error)

// carrierInitFuncs map contains the initilizer function for each carrier type.
var carrierInitFuncs map[string]carrierInitFunc

// RegisterCarrier should be used by the different carriers in order to register their `New` function
// Usually this method is called inside the `init` function of a carrier
func RegisterCarrier(carrierType string, initFunc carrierInitFunc) {
	if carrierInitFuncs == nil {
		carrierInitFuncs = make(map[string]carrierInitFunc)
	}
	carrierInitFuncs[carrierType] = initFunc
}

// New parse the carrier config and return a carrier object ready to be used.
func New(conf map[string]interface{}) (Carrier, error) {
	carrierType := conf["type"].(string)
	carrierInitializer, ok := carrierInitFuncs[carrierType]
	if !ok {
		return nil, fmt.Errorf("Unknown carrier type %q", carrierType)
	}

	carrier, err := carrierInitializer(conf)
	if err != nil {
		return nil, fmt.Errorf("Error initializing carrier %q: %v", carrier, err)
	}
	return carrier, nil
}
