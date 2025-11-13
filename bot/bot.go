package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot/botprom"
	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"
	"github.com/a-random-lemurian/go-owot-log/bot/chsize"
	"github.com/a-random-lemurian/go-owot-log/bot/cmdcore"
	"github.com/a-random-lemurian/go-owot-log/bot/lastseen"
	"github.com/a-random-lemurian/go-owot-log/bot/world"

	"github.com/a-random-lemurian/go-owot"
	
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Bot struct {
	startTime    time.Time
	v            *viper.Viper
	log          *zap.Logger
	db           Database
	cmdPrefix    string
	cmds         map[string]*chatcmd.Command
	cmdsMu       sync.Mutex
	msgCount     *chsize.MsgCounter
	trustedUsers map[string]any
	lastseen     *lastseen.Lastseen
	botprom *botprom.BotPrometheus
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
	b.log.Info("Initializing commands")
	b.cmdPrefix = b.v.GetString("command_prefix")

	msgCount := chsize.NewMsgCounter()
	msgCount.RegisterCommands(b)
	b.msgCount = msgCount

	lastseenService, err := lastseen.NewLastseen(b.v, b.log)
	if err != nil {
		return errors.Wrap(err, "initializing lastseen")
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

func (b *Bot) initPrometheus() error {
	if !b.v.IsSet("botprom.listen_addr") {
		return fmt.Errorf("botprom.listen_addr unset")
	}
	bp, err := botprom.NewBotPrometheus(b.v.GetString("botprom.listen_addr"), b.log)
	if err != nil {
		return errors.Wrap(err, "creating prometheus object")
	}
	b.botprom = bp
	return nil
}

func (b *Bot) Run() error {
	b.startTime = time.Now()

	log := GetZapLogger(b.v, "application.log")
	b.log = log.Named("bot")

	if err := b.getDatabase(); err != nil {
		return errors.Wrap(err, "initializing message database")
	}

	log.Info("Hi, good to see you")
	defer log.Info("Goodbye")

	wg := sync.WaitGroup{}

	if err := b.initCommands(); err != nil {
		return errors.Wrap(err, "initializing commands")
	}
	if err := b.initPrometheus(); err != nil {
		return errors.Wrap(err, "initializing prometheus")
	}

	b.initMiscConfig()
	b.initWorlds(&wg)

	wg.Add(1)
	go func () {
		defer wg.Done()
		err := b.botprom.Run()
		log.Error("prometheus died", zap.Error(err))
	}()

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

	b.log.Info(fmt.Sprintf("Registering new command %s, ratelimit 1 per %s", name, rl),
		zap.Duration("ratelimit", rl), zap.String("commandName", name))

	b.cmds[name] = cmd
}
