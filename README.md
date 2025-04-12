# go-owot-log

`go-owot-log` is a partial rewrite of Lemuria's owot_ws_log in the Go programming language.

The bot is very barebones and only writes messages to files, instead of using ClickHouse, out
of sheer necessity - Lemuria's VPS only had two gigabytes of RAM and could not handle such a
heavy database.
