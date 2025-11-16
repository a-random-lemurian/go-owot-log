# go-owot-log

`go-owot-log` is a partial rewrite of Lemuria's owot_ws_log in the Go programming language.

The bot is very barebones and only writes messages to files, instead of using ClickHouse, out
of sheer necessity - Lemuria's VPS only had two gigabytes of RAM and could not handle such a
heavy database.

Not to be confused with [go-owot](https://github.com/a-random-lemurian/go-owot), which
handles the OWOT websocket protocol itself.

## Run

The preferred method for running in production is to use Docker. Thankfully, it's as easy
as `make docker` and `docker compose up`.

To transfer to the VPS or whereever the production image is being run:
```
docker save lemuriaph/go-owot-log-main:latest | pv | ssh lemuria@vps1 docker load
```

(This command will only work for Lemuria; replace the ssh command with your own server)

## Versions

Version numbers aren't really necessary as this is rolling-release and code should be tested
and working before going to master, but when I switched to stdout logging with fluentd to
become 12-factor in commit `a83c81b783b66ced4ec3afb5435fbd81ac7d235e`, there would need to
be some transition work done.
