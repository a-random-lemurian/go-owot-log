package chsize

import (
	"fmt"
	"sync/atomic"

	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"
	"github.com/a-random-lemurian/go-owot-log/internal/util"

	"github.com/a-random-lemurian/go-owot"

	"github.com/dustin/go-humanize"
)

type MsgCounter struct {
	msgCount atomic.Int64
}

func NewMsgCounter() *MsgCounter {
	return &MsgCounter{}
}

func (m *MsgCounter) IngestMessage(c *owot.MessageChat) {
	m.msgCount.Add(1)
}

func (m *MsgCounter) ChSize(c chatcmd.Context, msg *owot.MessageChat) {
	size, err := util.DirectorySize(c.Viper().GetString("data_directory"))
	if err != nil {
		return
	}

	c.Chat(
		fmt.Sprintf(
			"data directory: %s / since bot start: %d msgs logged",
			humanize.Bytes(uint64(size)),
			m.msgCount.Load(),
		),
	)
}

func (m *MsgCounter) RegisterCommands(rc chatcmd.CommandRegistry) {
	rc.RegisterCommand("size", m.ChSize)
}
