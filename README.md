<div align="center">
    <h2>NotifyMe</h2>
    <p align="center">
        <p>Get notified when your command finished executing</p>
    </p>
</div>



## Contents

* [Installation](#installation)
* [Usage](#usage)
* [Notifiers](#notifiers)
    * [Messenger](#messenger)
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
### Prebuild binary

Prebuild binaries are available for both MacOS and Linux. They come in two flavors, 32 bits and 64 bits. Download links can be found in the table below:

|         | Linux           | MacOS  |
| ------------- |:-------------:|:-----:|
| 32 bits      | [Download](https://s3.amazonaws.com/clinotify.me/notifyme_linux32) | [Download](https://s3.amazonaws.com/clinotify.me/notifyme_darwin32) |
| 64 bits      | [Download](https://s3.amazonaws.com/clinotify.me/notifyme_linux64)      |   [Download](https://s3.amazonaws.com/clinotify.me/notifyme_darwin64) |

## Notifiers
Currently only Messenger is implemented as a notifier, others will be be implemented in the near future. Feel free to hack into the project  and implement other notifiers.
### Messenger
By using the Messenger notifier you will get notified via Facebook messenger.

First you need to get your `token` by talking to the [NotifyMe](https://www.facebook.com/clinotify.me/) Chat Bot. This can be done by sending `token` or `code` to the bot as shown in the picture below:
<p align="center">
    <img height=450 src=".github/MessengerCode.png">
</p>

Now that you have your token, edit your `~/.notifyme` config file and add it to `messenger_tokens` list:

```
{
    ...
    "messenger_tokens": [
        ...
        "YOUR_TOKEN_HERE"
    ],
    ...
}
```

**NOTE**: You can receive your token at any time by re-sending `token` or `code` to the Chat Bot.

## Usage

First you need to configure your [notifiers](#notifiers) by setting the right values (mainly tokens) in your `~/.notifyme` config file.

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
- etc.

is `NotifyMe` making your life simpler? tell us how are you using it :smile: !

## Todo

- Add flags through environment variables
- Enrich the configuration and add filters (e.g: only send erroned commands)
- Support other notifiers
    - Slack
- Add a re-execute functionality

## License

This repository has been released under the [MIT License](LICENSE)

------------------
Made with ♥ by [Think.iT](http://www.think-it.io/).