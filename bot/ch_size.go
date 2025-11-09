package bot

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/internal/util"
	"go.uber.org/zap"
)

// type CommandFunction func(*Bot,*World,*owot.MessageChat)

func chSize(b *Bot, w *World, c *owot.MessageChat) {
	size, err := util.DirectorySize(b.v.GetString("data_directory"))
	if err != nil {
		b.log.Error("oops", zap.Error(err))
	}
	w.Chat(fmt.Sprintf("data directory: %s", humanize.Bytes(uint64(size))), c.Location)
}
