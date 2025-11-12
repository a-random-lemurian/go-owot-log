package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"
	"github.com/a-random-lemurian/go-owot-log/bot/chsize"
	"github.com/a-random-lemurian/go-owot-log/bot/cmdcore"
	"github.com/a-random-lemurian/go-owot-log/bot/lastseen"
	"github.com/a-random-lemurian/go-owot-log/bot/world"

	"github.com/a-random-lemurian/go-owot"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Bot struct {
	startTime time.Time
	v         *viper.Viper
	log       *zap.Logger
	db        Database
	cmdPrefix string
	cmds      map[string]*chatcmd.Command
	cmdsMu    sync.Mutex
	msgCount  *chsize.MsgCounter
	trustedUsers map[string]any
	lastseen     *lastseen.Lastseen
}

func NewBot(v *viper.Viper) (*Bot, error) {
	return &Bot{v: v, cmds: make(map[string]*chatcmd.Command), trustedUsers: make(map[string]any)}, nil
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

func (b *Bot) initCommands() error {
	b.cmdPrefix = b.v.GetString("command_prefix")

	msgCount := chsize.NewMsgCounter()
	msgCount.RegisterCommands(b)
	b.msgCount = msgCount

	lastseenService, err := lastseen.NewLastseen(b.v)
	if err != nil {
		return err
	}
	b.lastseen = lastseenService
	b.lastseen.AttachCommands(b)

	b.RegisterCommand("uptime", cmdcore.Uptime)
	b.RegisterCommand("version", cmdcore.Version)
	b.RegisterCommand("help", cmdcore.Help)

	return nil
}

func (b *Bot) initWorlds(wg *sync.WaitGroup) {
	worlds := b.v.GetStringSlice("worlds")
	globalReceiver := b.v.GetString("global_receiver")

	for _, worldName := range worlds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.log.Info("Creating world", zap.String("worldName", worldName))
			acceptsGlobal := (worldName == globalReceiver)
			for {
				world, err := world.CreateWorld(worldName, b.log, b, acceptsGlobal)
				if err != nil {
					b.log.Error("Cannot create World", zap.String("worldName", worldName), zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}
				err = world.Run()
				if err != nil {
					b.log.Warn("Disconnected from world, reconnecting", zap.String("worldName", worldName), zap.Error(err))
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func (b *Bot) initMiscConfig() {
	b.v.SetDefault("chat.location.page.enabled", true)
	b.v.SetDefault("chat.location.global.enabled", true)

	trusted := b.v.GetStringSlice("users.trusted")
	for _, v := range trusted {
		b.trustedUsers[v] = 1
	}
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

	if err := b.initCommands(); err != nil {
		return err
	}

	b.initMiscConfig()
	b.initWorlds(&wg)

	wg.Wait()

	return nil
}

func (b *Bot) IsTrusted(chat *owot.MessageChat) bool {
	_, ok := b.trustedUsers[chat.RealUsername]
	return ok && b.v.GetBool("ratelimits.trusted_exempt")
}

func (b *Bot) RegisterCommand(name string, fn chatcmd.CommandFunction) {
	b.cmdsMu.Lock()
	defer b.cmdsMu.Unlock()

	// rlKey: ratelimit key
	rlKey := fmt.Sprintf("ratelimits.commands.%s", name)
	b.v.SetDefault(rlKey, 1*time.Second)
	rl := b.v.GetDuration(rlKey)

	cmd := &chatcmd.Command{
		Function:  fn,
		Ratelimit: rate.NewLimiter(rate.Every(rl), 1),
	}

	b.log.Info(fmt.Sprintf("Registering new command %s, ratelimit 1 per %s", name, rl), zap.Duration("ratelimit", rl))

	b.cmds[name] = cmd
}
