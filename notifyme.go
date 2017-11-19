package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"path/filepath"

	"github.com/kr/pty"
	"github.com/think-it-labs/clinotify.me/config"
	"github.com/think-it-labs/clinotify.me/notifier"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Setup the logger
	debugModeEnabled := strings.ToLower(os.Getenv("DEBUG")) == "true"
	if debugModeEnabled {
		log.SetLevel(log.DebugLevel)
	}

}

func main() {

	// Get the current user dir and parse the config file
	user, _ := user.Current()
	configPath := filepath.Join(user.HomeDir, ".notifyme")
	config, err := config.FromFile(configPath)
	if err != nil {
		exitErrConfig()
	}

	userCmd := os.Args[1:]
	if len(userCmd) == 0 {
		exitUsage()
	}

	shell := os.Getenv("SHELL")
	args := []string{
		"-i",
		"-c",
		strings.Join(userCmd, " "),
	}
	cmd := exec.Command(shell, args...)

	argsLog := sliceString(args)
	log.Debugf("Command: %s %s", shell, argsLog)

	var exitCode int

	f, _ := pty.Start(cmd)

	// Create a buffer and add it to a multiwriter
	var output bytes.Buffer
	dataWriter := io.MultiWriter(os.Stdout, &output)

	// Copy tty output to our multiwriter
	io.Copy(dataWriter, f)

	// Wait for the process to exit and get it's status code
	err = cmd.Wait()
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = 1 // default exit code
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	// Build the notification
	notification := notifier.Notification{
		Cmd:      strings.Join(userCmd, " "),
		ExitCode: exitCode,
		Logs:     output.Bytes(),
	}

	// Make notifiers and send notification to them
	// TODO: make this dynamic.
	var notifiers []notifier.Notifier
	notifiers = append(notifiers, notifier.Messenger{
		Token:        config.MessengerToken,
		Notification: notification,
	})

	for _, notifier := range notifiers {
		err := notifier.Notify()
		if err != nil {
			log.Errorf("Error sending notification: %s", err)
		}
	}

	// Exit and use the same user's command exitCode
	os.Exit(exitCode)
}

func exitErrConfig() {
	fmt.Println(`Your config file is missing or something is wrong with it.
Please make sure you have the file ~/.notifyme and that it's content looks like:

{
    "messenger_token": "YOUR_MESSENGER_TOKEN_HERE",
}`)
	os.Exit(1)
}

func exitUsage() {
	fmt.Printf("%s CMD_HERE ARG1 ARG2 ...\n", os.Args[0])
	os.Exit(2)
}

func sliceString(slice []string) string {
	sliceStr := fmt.Sprintf("%s", slice)[1:]
	sliceStr = sliceStr[:len(sliceStr)-1]
	return sliceStr
}
