package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/think-it-labs/notifyme/carriers"
	"github.com/think-it-labs/notifyme/notification"

	isatty "github.com/mattn/go-isatty"
	"github.com/think-it-labs/notifyme/argparser"
	"github.com/think-it-labs/notifyme/command"
	"github.com/think-it-labs/notifyme/config"

	// import carriers so they register they initialize function
	_ "github.com/think-it-labs/notifyme/carriers/file"
	_ "github.com/think-it-labs/notifyme/carriers/messenger"
	_ "github.com/think-it-labs/notifyme/carriers/slack"

	log "github.com/sirupsen/logrus"
)

var arguments argparser.Arguments

func init() {
	arguments = argparser.MustParse(os.Args[1:])

	// Set default config path if configfile is empty
	if arguments.ConfigFile == "" {
		arguments.ConfigFile = config.DefaultConfigPath
	}

	// Setup log level
	if len(arguments.Verbose) == 0 {
		log.SetLevel(log.WarnLevel)
	} else if len(arguments.Verbose) == 2 {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	// Parse the config file
	log.Infof("Config file: %s", arguments.ConfigFile)
	cfg, err := config.FromFile(arguments.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) && arguments.ConfigFile == config.DefaultConfigPath {
			exitNoDefaultConfig()
		} else {
			exitBadConfig(err)
		}
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		log.Warnln("It seems like the output is piped, please refer to https://clinotify.me/piped for more info about this.")
	}

	cmd := command.New(arguments.UserCmd)
	log.Infof("Command: %s", strings.Join(arguments.UserCmd, " "))

	// Setup stdout and stderr writers
	output := new(bytes.Buffer)
	cmd.AddStdoutWriter(output)
	cmd.AddStderrWriter(output)

	// Start the command and wait for it
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Cannot start the command: %s\n", err)
	}
	exitCode := cmd.Wait()

	// Build the notification
	notif := notification.New(strings.Join(arguments.UserCmd, " "),
		exitCode,
		output.Bytes(),
	)

	// Build carriers
	var carriersList []carriers.Carrier
	for _, carrierConf := range cfg.Carriers {
		carrier, err := carriers.New(carrierConf)
		if err != nil {
			log.Errorf("Error initializing carrier: %v", err)
			continue
		}
		carriersList = append(carriersList, carrier)
	}

	log.Infof("%d carrier ready", len(carriersList))

	// Send notifications
	var wg sync.WaitGroup
	wg.Add(len(carriersList))
	for _, carrier := range carriersList {
		go func(carrier carriers.Carrier) {
			err := carrier.Send(notif)
			if err != nil {
				log.Errorf("Error sending notification: %v", err)
			}
			wg.Done()
		}(carrier)

	}
	wg.Wait()

	// Exit and use the same user's command exitCode
	os.Exit(exitCode)
}

func exitBadConfig(err error) {
	configFile := arguments.ConfigFile
	fmt.Fprintf(os.Stderr, "Error reading config from %s: %v\n", configFile, err)
	os.Exit(1)
}

func exitNoDefaultConfig() {
	configFile := config.DefaultConfigPath
	fmt.Printf("Configuration file %s cannot be found. Making one for you :)\n", configFile)
	err := config.CreateDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create configuration file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration file %s have been created, please edit it to start using NotifyMe\n", configFile)
	os.Exit(0)
}
