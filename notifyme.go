package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	isatty "github.com/mattn/go-isatty"
	"github.com/think-it-labs/notifyme/argparser"
	"github.com/think-it-labs/notifyme/command"
	"github.com/think-it-labs/notifyme/config"
	"github.com/think-it-labs/notifyme/notification"

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
	notificationData := notification.NotificationData{
		Cmd:      strings.Join(arguments.UserCmd, " "),
		ExitCode: exitCode,
		Logs:     output.Bytes(),
	}

	// Build the list of notification to be sent
	var notifications []notification.Notification
	if cfg.MessengerEnabled {
		for _, token := range cfg.MessengerTokens {
			if token == "" {
				continue
			}
			notifications = append(notifications, notification.Messenger{
				Token:            token,
				NotificationData: notificationData,
			})
		}
	}

	log.Infof("Sending %d notification(s)", len(notifications))

	// Send notifications
	var wg sync.WaitGroup
	wg.Add(len(notifications))
	for _, notif := range notifications {
		go func(n notification.Notification) {
			err := n.Send()
			if err != nil {
				log.Errorf("Error sending notification: %s", err)
			}
			wg.Done()
		}(notif)

	}
	wg.Wait()

	// Exit and use the same user's command exitCode
	os.Exit(exitCode)
}

func exitBadConfig(err error) {
	configFile := "BLABL"
	fmt.Fprintf(os.Stderr, "Error reading config from %s: %s\n", configFile, err)
	os.Exit(1)
}

func exitNoDefaultConfig() {
	configFile := config.DefaultConfigPath
	fmt.Printf("Configuration file %s cannot be found. Making one for you :)\n", configFile)
	err := config.CreateDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create configuration file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration file %s have been created, please edit it to start using NotifyMe\n", configFile)
	os.Exit(0)
}
