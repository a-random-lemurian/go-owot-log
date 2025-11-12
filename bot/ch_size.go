package bot

import (
	"fmt"

	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/bot/world"
	"github.com/a-random-lemurian/go-owot-log/internal/util"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func chSize(b *Bot, w *world.World, c *owot.MessageChat) {
	size, err := util.DirectorySize(b.v.GetString("data_directory"))
	if err != nil {
		b.log.Error("oops", zap.Error(err))
	}
	w.Chat(
		fmt.Sprintf(
			"data directory: %s / since bot start: %d msgs logged",
			humanize.Bytes(uint64(size)),
			b.msgCount.Load(),
		),
		c.Location,
	)
}
