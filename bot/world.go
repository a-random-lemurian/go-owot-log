package bot

import (
	"encoding/json"
	"time"

	"github.com/a-random-lemurian/go-owot"
	"github.com/a-random-lemurian/go-owot-log/internal/types"
	"go.uber.org/zap"
)

type World struct {
	owotConn *owot.OwotConn
	url      string
	database Database
	appLog   *zap.Logger

	AcceptGlobal bool
}

// TODO allow passing in sites, also allow multi-site logging?
// Who knows if Fern will set up her own instance of NWOT, a
// permanent one?

func CreateOWOTUrl(worldName string) string {
	// Detect the worldName being ""
	url := "wss://ourworldoftext.com/" + worldName
	if worldName != "" {
		url += "/"
	}
	url += "ws/"
	return url
}

func CreateWorld(worldName string, log *zap.Logger, db Database, acceptGlobal bool) (*World, error) {
	url := CreateOWOTUrl(worldName)
	conn, err := owot.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	world := &World{
		owotConn:     conn,
		database:     db,
		url:          url,
		appLog:       log,
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
	time := time.Now()

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

	err = w.database.AddMessage(&types.Message{Origin: w.url, WebsocketMessage: data, Time: time})
	if err != nil {
		w.appLog.Error("Failed to insert message into database", zap.Error(err), zap.ByteString("rawMessage", data))
	}
}

func (w *World) Run() error {
	return w.owotConn.Run()
}
