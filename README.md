CLI client that connects to TP-Link Archer MR600

There are 3 implemented commands that handle SMS:

- List Inbox
- Send SMS message
- Find SMS match by content and sender

## How to Run

Go is required to run, after installing Go run:

```sh
go mod download
```

Then you can run with:

```sh
go run main.go sms list
```

Or you can install as go command with:

```sh
go install

tp-link sms list
```

### Config

Config is saved in `$HOME/.tp-link.yaml`, it looks like and defaults to:

```yaml
host: http://192.168.1.1
password: admin
username: admin
```

### Script

You can create a script like:

```sh
#!/bin/sh
if tp-link sms find --phrase=some-phrase --from=sender; then
	tp-link sms send --to=destination --content=some-content
fi
```

To check inbox for received sms and send sms back if found.
