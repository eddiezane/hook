# Captain Hook

Captain Hook is a CLI tool for firing a collection of known webhooks for development.

Have you ever needed to build an app that consumes a new [pull request webhook event](https://developer.github.com/v3/activity/events/types/#pullrequestevent) from GitHub? Or maybe using Twilio's [inbound SMS webhook](https://www.twilio.com/docs/usage/webhooks/sms-webhooks). Sure there's some tooling that lets you replay those events once they've been made, but what if you don't want to jump through all of the hoops of creating a test repo, configuring the webhook to your dev server, and opening a new pull request.

That's the friction Captain Hook aims to solve.

## Implmentation

Webhooks are serialized to and from a basic YAML syntax with the intention of being human creatable and editable.

# Installation 

Until there are binaries built you can install Captain Hook into your `$GOPATH/bin`

```
go get -u github.com/eddiezane/captain-hook
```

# Usage

Still a work in progress but a general idea of firing a webhook is:

```bash
hook fire webhooks/twilio/sms http://localhost:8080
```

Captain Hook also has an HTTP server for recording new webhooks:

```bash
hook record --port 8080 path/to/new/webhook.yml
```

# Roadmap

- [x] Basic working POC
  - [x] Fire command
  - [x] Record command
- [ ] Initial release candidate
  - [ ] Basic collection of webhooks to convey usability (Twilio, GitHub, ...)
  - [ ] Don't use default http client
  - [ ] Server error handling
  - [ ] Server shutdown logic
  - [ ] Better error handling in current commands
  - [ ] Implement proper flags
  - [ ] Add view command to view a webhook in it's YAML format
- [ ] Catalog logic
  - [ ] Define spec for a catalog
  - [ ] Download and lookup (tap) a new catalog
  - [ ] Create default catalog as it's own GitHub repo
- [ ] Template logic for webhooks (sub in vars)
- [ ] Web UI
- [ ] CI/CD logic for consuming webhooks

# License

MIT
