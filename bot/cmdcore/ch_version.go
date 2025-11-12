package cmdcore

import (
	"fmt"

	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"
	"github.com/carlmjohnson/versioninfo"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func Version(chatCtx chatcmd.Context, c *owot.MessageChat) {
	chatCtx.Chat(fmt.Sprintf("go-owot-log version %s", versioninfo.Short()))
}
