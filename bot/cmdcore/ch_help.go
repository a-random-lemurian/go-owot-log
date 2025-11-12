package cmdcore

import (
	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"

	"github.com/a-random-lemurian/go-owot"
)

func Help(chatCtx chatcmd.Context, c *owot.MessageChat) {
	// TODO: maybe require commands to include help messages when registering themselves
	chatCtx.Chat("go-owot-log is a bot by Lemuria, the successor to owot_ws_log! / ch size - size of chatlogs / other cmds: help, uptime, version, lastseen, lastseen-consent / to opt out of lastseen: 'ch lastseen-consent no'")
}
