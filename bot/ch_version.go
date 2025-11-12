package bot

import (
	"fmt"

	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/bot/world"
	"github.com/carlmjohnson/versioninfo"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func chVersion(b *Bot, w *world.World, c *owot.MessageChat) {
	w.Chat(fmt.Sprintf("go-owot-log version %s", versioninfo.Short()), c.Location)
}
