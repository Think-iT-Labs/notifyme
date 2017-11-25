package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
	"sync"

	"path/filepath"

	"github.com/think-it-labs/notifyme/command"
	"github.com/think-it-labs/notifyme/config"
	"github.com/think-it-labs/notifyme/notification"

	log "github.com/sirupsen/logrus"
)

var configFile string

func init() {
	// Setup the logger
	debugModeEnabled := strings.ToLower(os.Getenv("DEBUG")) == "true"
	if debugModeEnabled {
		log.SetLevel(log.DebugLevel)
	}

	// Config file path
	user, _ := user.Current()
	configFile = filepath.Join(user.HomeDir, ".notifyme")
	if os.Getenv("NOTIFYME_CONFIG_FILE") != "" {
		configFile = os.Getenv("NOTIFYME_CONFIG_FILE")
	}

}

func main() {

	config, err := config.FromFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			exitNoConfig()
		} else {
			exitBadConfig(err)
		}
	}

	userCmd := os.Args[1:]
	if len(userCmd) == 0 {
		exitUsage()
	}

	cmd := command.New(userCmd)
	log.Debugf("Command: %s", strings.Join(userCmd, " "))

	tty, err := cmd.StartWithTTY()
	if err != nil {
		log.Fatalf("Cannot start the command with a new tty: %s\n", err)
	}
	defer tty.Close()

	go func() {
		io.Copy(tty, os.Stdin)
	}()

	// Create a buffer and add it to a multiwriter
	var output bytes.Buffer
	dataWriter := io.MultiWriter(os.Stdout, &output)

	// Copy tty output to our multiwriter
	io.Copy(dataWriter, tty)

	// Wait for the process to exit and get it's status code
	exitCode := cmd.Wait()

	// Build the notification
	notificationData := notification.NotificationData{
		Cmd:      strings.Join(userCmd, " "),
		ExitCode: exitCode,
		Logs:     output.Bytes(),
	}

	// Build the list of notification to be sent
	var notifications []notification.Notification
	if config.MessengerEnabled {
		for _, token := range config.MessengerTokens {
			if token == "" {
				continue
			}
			notifications = append(notifications, notification.Messenger{
				Token:            token,
				NotificationData: notificationData,
			})
		}
	}

	log.Debugf("Sending %d notification(s)", len(notifications))

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
	fmt.Fprintf(os.Stderr, "Error reading config from %s: %s\n", configFile, err)
	os.Exit(1)
}

func exitNoConfig() {
	fmt.Printf("Configuration file %s cannot be found. Making one for you :)\n", configFile)
	err := config.CreateDefault(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create configuration file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration file %s have been created, please edit it to start using notify.me\n", configFile)
	os.Exit(0)
}

func exitUsage() {
	fmt.Fprintf(os.Stderr, "%s CMD_HERE ARG1 ARG2 ...\n", os.Args[0])
	os.Exit(2)
}

func sliceString(slice []string) string {
	sliceStr := fmt.Sprintf("%s", slice)[1:]
	sliceStr = sliceStr[:len(sliceStr)-1]
	return sliceStr
}
