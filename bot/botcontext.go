package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/world"

	"github.com/a-random-lemurian/go-owot"

	"github.com/spf13/viper"
)

type botContext struct {
	b         *Bot
	w         *world.World
	c         *owot.MessageChat
	startTime time.Time
}

func (bc *botContext) Chat(message string) {
	final := ""
	if bc.c.PrivateMessage != nil && *bc.c.PrivateMessage == owot.PMTToMe {
		final += fmt.Sprintf("/tell %s ", bc.c.Id)
	}
	final += message
	bc.w.Chat(final, bc.c.Location.String())
}

func (bc *botContext) IsTrusted(username string) bool {
	// probably should refactor *Bot.IsTrusted ....
	return bc.b.IsTrusted(&owot.MessageChat{RealUsername: username})
}

func (bc *botContext) Viper() *viper.Viper {
	return bc.b.v
}

func (bc *botContext) StartTime() time.Time {
	return bc.startTime
}

func (b *Bot) handleCommand(w *world.World, chat *owot.MessageChat) {
	fields := strings.Fields(chat.Message)
	if len(fields) < 2 {
		return
	}

	chatCtx := &botContext{
		b:         b,
		w:         w,
		c:         chat,
		startTime: b.startTime,
	}

	commandName := fields[1]
	command, ok := b.cmds[commandName]
	if ok {
		if command.Ratelimit.AllowN(time.Now(), 1) || b.IsTrusted(chat) {
			command.Function(chatCtx, chat)
		}
	}
}
