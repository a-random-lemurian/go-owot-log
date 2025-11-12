package world

import (
	"fmt"
	"net/url"
	"time"

	"github.com/a-random-lemurian/go-owot"
	"github.com/carlmjohnson/versioninfo"
	"go.uber.org/zap"
)

type RawSink interface {
	HandleRaw(*World, []byte)
}

type World struct {
	owotConn *owot.OwotConn
	url      string
	appLog   *zap.Logger
	nickname string
	color    string
	rawSink  RawSink

	AcceptGlobal bool
}

// TODO allow passing in sites, also allow multi-site logging?
// Who knows if Fern will set up her own instance of NWOT, a
// permanent one?

func CreateOWOTUrl(worldName string, hide bool) string {
	url := &url.URL{
		Scheme: "wss",
		Host:   "ourworldoftext.com",
		Path:   (worldName + "/") + "ws/",
	}

	if hide {
		q := url.Query()
		q.Add("hide", "1")
		url.RawQuery = q.Encode()
	}

	return url.String()
}

func CreateWorld(worldName string, log *zap.Logger, rawSink RawSink, acceptGlobal bool) (*World, error) {
	url := CreateOWOTUrl(worldName, true)
	conn, err := owot.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	world := &World{
		owotConn:     conn,
		url:          url,
		appLog:       log,
		AcceptGlobal: acceptGlobal,
		rawSink: rawSink,
	}

	conn.HandleRaw = world.handleRaw

	world.appLog.Info("Created world",
		zap.String("worldName", worldName),
		zap.String("url", url),
		zap.Bool("acceptsGlobal", acceptGlobal),
	)

	return world, nil
}

func (w *World) handleRaw(data []byte) {
	w.rawSink.HandleRaw(w, data)
}

func (w *World) GetURL() string {
	return w.url
}

func (w *World) Run() error {
	return w.owotConn.Run()
}

// A wrapper for owotConn chat.
func (w *World) Chat(message string, location string) error {
	msgMeta := map[string]any{}
	msgMeta["version"] = fmt.Sprintf("go-owot-log version %s", versioninfo.Short())
	msgMeta["sent-date"] = time.Now().UTC().UnixMilli()
	return w.owotConn.Chat(w.nickname, message, location, w.color, msgMeta)
}
