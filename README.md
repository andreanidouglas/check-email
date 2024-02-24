# Check Email

This is a simple tool to check how many unread emails you have

Create a file `email.yml` (see example) to configure email server and credentials

## Run

```bash
$ ./check_email email.yml
```

## Building

_go minimum version: 1.20_

```bash
go build -o check_status cmd/main.go
```


## i3Status

**If you want this to execute on your i3status bar.**

Replace the i3status for the script `i3_common.sh`

```bash
bar {
    status_command i3_command.sh
}
```

