package bot

import (
	"net/url"

	"github.com/a-random-lemurian/go-owot"
	"go.uber.org/zap"
)

type World struct {
	owotConn *owot.OwotConn
	url      string
	appLog   *zap.Logger
	bot *Bot

	AcceptGlobal bool
}

// TODO allow passing in sites, also allow multi-site logging?
// Who knows if Fern will set up her own instance of NWOT, a
// permanent one?

func CreateOWOTUrl(worldName string, hide bool) string {
	url := &url.URL{
		Scheme: "wss",
		Host: "ourworldoftext.com",
		Path: (worldName + "/") + "ws/",
	}

	if hide {
		q := url.Query()
		q.Add("hide", "1")
		url.RawQuery = q.Encode()
	}

	return url.String()
}

func (b *Bot) CreateWorld(worldName string, log *zap.Logger, acceptGlobal bool) (*World, error) {
	url := CreateOWOTUrl(worldName, true)
	conn, err := owot.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	world := &World{
		owotConn:     conn,
		url:          url,
		appLog:       log,
		bot: b,
		AcceptGlobal: acceptGlobal,
	}

	conn.HandleRaw = world.HandleRaw

	world.appLog.Info("Created world",
		zap.String("worldName", worldName),
		zap.String("url", url),
		zap.Bool("acceptsGlobal", acceptGlobal),
	)

	return world, nil
}

func (w *World) HandleRaw(data []byte) {
	w.bot.handleRaw(w, data)
}

func (w *World) Run() error {
	return w.owotConn.Run()
}
