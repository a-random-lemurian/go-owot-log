package cmdcore

import (
	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"
)

func Help(chatCtx chatcmd.Context, c *owot.MessageChat) {
	chatCtx.Chat("go-owot-log is a bot by Lemuria, the successor to owot_ws_log! / ch size - size of chatlogs / other cmds: help, uptime, version")
}
