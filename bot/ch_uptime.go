package bot

import (
	"fmt"
	"time"

	"github.com/a-random-lemurian/go-owot"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func chUptime(b *Bot, w *World, c *owot.MessageChat) {
	startTimeString := b.startTime.UTC().Format("2006-01-02 15:04:05.000") + " UTC"
	w.Chat(fmt.Sprintf("go-owot-log uptime %s - start time %s", time.Since(b.startTime.UTC()), startTimeString), c.Location)
}
