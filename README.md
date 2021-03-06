# Hook

Hook is a CLI tool for firing a collection of known webhooks for development.

Have you ever needed to build an app that consumes a new [pull request webhook event](https://developer.github.com/v3/activity/events/types/#pullrequestevent) from GitHub? Or maybe using Twilio's [inbound SMS webhook](https://www.twilio.com/docs/usage/webhooks/sms-webhooks). Sure there's some tooling that lets you replay those events once they've been made, but what if you don't want to jump through all of the hoops of creating a test repo, configuring the webhook to your dev server, and opening a new pull request.

That's the friction Hook aims to solve.

## Implementation

Webhooks are serialized to and from a basic YAML syntax with the intention of being human creatable and editable.

# Installation

```bash
go get -u github.com/eddiezane/hook
```

# Usage

## Fire

Hooks can be fired locally by specifying the path:

```bash
hook fire webhooks/twilio/sms http://localhost:8080
```

File suffixes are fuzzy matched - specifying a hook file `foo` will match `foo`, `foo.yaml`, or `foo.yml`

### Catalogs

`hook` can be configured to read from remote Git repositories for hook data.

By default, hook comes installed with a default catalog of contributed hooks stored at https://github.com/eddiezane/hook-catalog.

```bash
hook fire @github/push http://localhost:8080
```

Additional catalogs can be configured via the `hook catalog` subcommand.

## Record

Hook also has an HTTP server for recording new webhooks:

```bash
hook record --port 8080 path/to/new/webhook.yml
```

Multiple hooks received by the server will be stored in the same file as a multidoc yaml (separated by `---`).

# Roadmap

- [x] Basic working POC
  - [x] Fire command
  - [x] Record command
- [ ] Initial release candidate
  - [ ] Basic collection of webhooks to convey usability (Twilio, GitHub, ...)
  - [ ] Don't use default http client
  - [ ] Server error handling
  - [ ] Server shutdown logic
  - [x] Better error handling in current commands
  - [x] Implement proper flags
  - [ ] Add view command to view a webhook in it's YAML format
- [x] Catalog logic
  - [x] Define spec for a catalog
  - [x] Download and lookup (tap) a new catalog
  - [x] Create default catalog as it's own GitHub repo
  - [ ] Add automatic workflows to update webhooks.
- [ ] Template logic for webhooks (sub in vars)
- [ ] Web UI

# License

MIT
