package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/world"

	"github.com/a-random-lemurian/go-owot-log/internal/types"

	"github.com/a-random-lemurian/go-owot"

	"go.uber.org/zap"
)

func (b *Bot) HandleRaw(w *world.World, data []byte) {
	type partial struct {
		Kind     string `json:"kind"`
		Location string `json:"location"`
	}
	var kind partial
	err := json.Unmarshal(data, &kind)
	if err != nil {
		b.log.Error("Failed to parse incoming message", zap.Error(err), zap.ByteString("rawMessage", data))
		return
	}

	// TODO more robust blacklist system
	if kind.Kind == "tileUpdate" {
		return
	}
	if kind.Kind == "chat" && kind.Location == "global" {
		if !w.AcceptGlobal {
			return
		}
	}

	b.handleChat(w, data)
}

func (b *Bot) handleChat(w *world.World, data []byte) {
	time := time.Now()

	var chat owot.MessageChat
	if err := json.Unmarshal(data, &chat); err != nil {
		b.log.Error("failed to parse message", zap.Error(err), zap.ByteString("rawData", data))
	}

	b.lastseen.IngestMessage(&chat)
	b.msgCount.IngestMessage(&chat)

	shouldReply := b.v.GetBool(fmt.Sprintf("chat.location.%s.enabled", chat.Location))
	if strings.HasPrefix(chat.Message, b.cmdPrefix) && shouldReply {
		b.handleCommand(w, &chat)
	}

	err := b.db.AddMessage(&types.Message{Origin: w.GetURL(), WebsocketMessage: data, Time: time})
	if err != nil {
		b.log.Error("Failed to insert message into database", zap.Error(err), zap.ByteString("rawMessage", data))
	}
}
