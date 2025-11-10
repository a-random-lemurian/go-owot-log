package bot

import (
	"github.com/a-random-lemurian/go-owot"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func chHelp(b *Bot, w *World, c *owot.MessageChat) {
	w.Chat("go-owot-log is a bot by Lemuria, the successor to owot_ws_log! / ch size - size of chatlogs / other cmds: help, uptime, version", c.Location)
}
