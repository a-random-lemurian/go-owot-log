package cmdcore

import (
	"fmt"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"

	"github.com/a-random-lemurian/go-owot"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func Uptime(chatCtx chatcmd.Context, c *owot.MessageChat) {
	startTimeString := chatCtx.StartTime().UTC().Format("2006-01-02 15:04:05.000") + " UTC"
	chatCtx.Chat(fmt.Sprintf("go-owot-log uptime %s - start time %s", time.Since(chatCtx.StartTime().UTC()), startTimeString))
}
