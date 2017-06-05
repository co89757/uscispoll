# USCIS Case Progress Monitor in Golang
A simple USCIS case poller written in Go to poll the case status and provide notification service on status change. Once set up on a server, it then starts polling the status of your case daily and send email notice on change.

## How to install
if you have go installed in your environemnt, then use go get:
```
go get github.com/co89757/uscispoll
```

### Pre-built binary
Prebuilt binary is available via

## Usage
It's a command-line utility that takes a single argument: your receipt number (aka. case number) and a configuration regarding email notice setup
### Configure Email settings
All the email related settings are in `emailcfg.json` file, populate it with your own credentials for email sending.
### Example usage
After you fetch it with go get, the binary is in your `$GOPATH/bin` directory
```
$GOPATH/bin/uscispoll -case <case_number>
```
