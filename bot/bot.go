package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/internal/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type CommandFunction func(*Bot, *World, *owot.MessageChat)

type Command struct {
	Function  CommandFunction
	Ratelimit *rate.Limiter
}

type Bot struct {
	startTime time.Time
	v         *viper.Viper
	log       *zap.Logger
	db        Database
	cmdPrefix string
	cmds      map[string]*Command
	cmdsMu    sync.Mutex
	msgCount  atomic.Int64
	trustedUsers map[string]any
}

func NewBot(v *viper.Viper) (*Bot, error) {
	return &Bot{v: v, cmds: make(map[string]*Command), trustedUsers: make(map[string]any)}, nil
}

func (b *Bot) registerCommand(name string, fn CommandFunction) {
	b.cmdsMu.Lock()
	defer b.cmdsMu.Unlock()

	// rlKey: ratelimit key
	rlKey := fmt.Sprintf("ratelimits.commands.%s", name)
	b.v.SetDefault(rlKey, 1*time.Second)
	rl := b.v.GetDuration(rlKey)

	cmd := &Command{
		Function:  fn,
		Ratelimit: rate.NewLimiter(rate.Every(rl), 1),
	}

	b.log.Info(fmt.Sprintf("Registering new command %s, ratelimit 1 per %s", name, rl), zap.Duration("ratelimit", rl))

	b.cmds[name] = cmd
}

func (b *Bot) getDatabase() error {
	database, err := CreatePlainTextFileDatabase(b.v)
	if err != nil {
		b.log.Error("Cannot create PlainTextFileDatabase", zap.Error(err))
		return err
	}
	b.db = database
	return nil
}

func (b *Bot) Run() error {
	b.startTime = time.Now()
	log := GetZapLogger(b.v, "application.log")
	b.log = log

	if err := b.getDatabase(); err != nil {
		return err
	}

	log.Info("Hi, good to see you")
	defer log.Info("Goodbye")

	wg := sync.WaitGroup{}

	b.cmdPrefix = b.v.GetString("command_prefix")

	trusted := b.v.GetStringSlice("users.trusted")
	for _, v := range trusted {
		b.trustedUsers[v] = 1
	}

	b.registerCommand("size", chSize)
	b.registerCommand("uptime", chUptime)
	b.registerCommand("version", chVersion)
	b.registerCommand("help", chHelp)

	worlds := b.v.GetStringSlice("worlds")
	globalReceiver := b.v.GetString("global_receiver")

	for _, worldName := range worlds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Creating world", zap.String("worldName", worldName))
			acceptsGlobal := (worldName == globalReceiver)
			for {
				world, err := b.CreateWorld(worldName, log, acceptsGlobal)
				if err != nil {
					log.Error("Cannot create World", zap.String("worldName", worldName), zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}
				err = world.Run()
				if err != nil {
					log.Warn("Disconnected from world, reconnecting", zap.String("worldName", worldName), zap.Error(err))
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	wg.Wait()

	return nil
}

func (b *Bot) handleRaw(w *World, data []byte) {
	type partial struct {
		Kind     string `json:"kind"`
		Location string `json:"location"`
	}
	var kind partial
	err := json.Unmarshal(data, &kind)
	if err != nil {
		w.appLog.Error("Failed to parse incoming message", zap.Error(err), zap.ByteString("rawMessage", data))
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

func (b *Bot) handleChat(w *World, data []byte) {
	time := time.Now()

	var chat owot.MessageChat
	if err := json.Unmarshal(data, &chat); err != nil {
		b.log.Error("failed to parse message", zap.Error(err), zap.ByteString("rawData", data))
	}

	b.msgCount.Add(1)

	if strings.HasPrefix(chat.Message, b.cmdPrefix) {
		b.handleCommand(w, &chat)
	}

	err := b.db.AddMessage(&types.Message{Origin: w.url, WebsocketMessage: data, Time: time})
	if err != nil {
		w.appLog.Error("Failed to insert message into database", zap.Error(err), zap.ByteString("rawMessage", data))
	}
}

func (b *Bot) isTrusted(chat *owot.MessageChat) bool {
	_, ok := b.trustedUsers[chat.RealUsername]
	return ok && b.v.GetBool("ratelimits.trusted_exempt")
}

func (b *Bot) handleCommand(w *World, chat *owot.MessageChat) {
	fields := strings.Fields(chat.Message)
	if len(fields) < 2 {
		return
	}

	commandName := fields[1]
	command, ok := b.cmds[commandName]
	if ok {
		if command.Ratelimit.AllowN(time.Now(), 1) || b.isTrusted(chat) {
			command.Function(b, w, chat)
		}
	}
}
