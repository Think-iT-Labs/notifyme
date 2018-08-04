<div align="center">
    <h2>NotifyMe</h2>
    <p align="center">
        <p>Get notified when your command finished executing</p>
    </p>
</div>



## Contents

* [Installation](#installation)
* [Usage](#usage)
* [Carriers](#carriers)
    * [Slack](#slack)
* [Use cases](#use-cases)
* [Todo](#todo)
* [License](#license)

## Installation

### Build from the source

`NotifyMe` client is written in Golang, to build it from the source you need to have `go` installed and your `GOPATH` configured (default to `~/go` in go 1.9)

Once done, `get` the code by running:
```shell
go get github.com/think-it-labs/notifyme
```

**NOTE**: the command above will download the repo into your `GOPATH` and build it for you. The resulting binary can be found in `$GOPATH/bin`, we recommend adding `$GOPATH/bin` to your `$PATH`.

## Carriers

A carrier is a plugin that can deliver a notification. 
Currently only Slack is implemented as a carrier, others will be be implemented in the near future. Feel free to hack into the project and implement other notifiers.
### Slack
By using the Slack carrier you will get notified via slack.

First you need to get your `token` by visiting this page [Slack Token](https://api.slack.com/custom-integrations/legacy-tokens).

Now that you have your token, edit your `~/.notifyme` config file and add the slack carrier. 

Example:
```
carriers:
  - type: slack
    token: "xoxp-XXXXXX"
    channels:
    - "@user1"
    - "#general"
```

## Usage

First you need to configure your [carriers](#carriers) by setting the right values (mainly tokens) in your `~/.notifyme` config file.

Second prepend `notifyme` to your command to get notified when it is finished executing.
```
$ notifyme COMMAND ARG1 ARG2 ...
```

For example to get notified when a Make build is finished, the command will look like:

```
$ notifyme make -j 4
```
## Use Cases
Use cases for `NotifyMe` are numerous, and here are some tasks that developers regulary  want to receive status notifications from.

- Cron jobs,
- Long running builds,
- Backups,
- Result of long running commands,
- etc.

Is `NotifyMe` making your life simpler? tell us how are you using it :smile: !

## Todo

[ ] Add flags through environment variables
[ ] Enrich the configuration and add filters (e.g: only send erroned commands)
[ ] Support other notifiers
    [X] Slack
    [X] Email
    [ ] File

## License

This repository has been released under the [MIT License](LICENSE)

------------------
Made with ♥ by [Think.iT](http://www.think-it.io/).