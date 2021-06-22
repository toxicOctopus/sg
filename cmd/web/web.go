package main

import (
	"context"
	"flag"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/toxicOctopus/sg/internal/centrifugo"
	"github.com/toxicOctopus/sg/internal/config"
	"github.com/toxicOctopus/sg/internal/database"
	"github.com/toxicOctopus/sg/internal/pkg/game"
	"github.com/toxicOctopus/sg/internal/twitch"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	jsClient     string
	startTime    time.Time
	env          config.Env
	globalConfig config.LiveConfig
)

func main() {
	ctx := context.Background()

	logrus.Info(ctx, "Starting up @ "+startTime.String())

	cfg := globalConfig.GetCfg()
	pgConfig := cfg.Postgres
	db, err := database.GetDB(ctx, pgConfig.Host, pgConfig.Port, pgConfig.Scheme, pgConfig.User, pgConfig.Password)
	defer db.Close(ctx)
	if err != nil {
		logrus.Fatalf("DB connection error: %s", err)
	}

	//TODO listen postgres changes for rc
	registeredChannels, err := database.GetRegisteredChannels(ctx, db)
	if err != nil {
		logrus.Fatal(ctx, err)
	}

	centrifugoClient, err := centrifugo.GetClient(cfg.Centrifugo.URL, cfg.Centrifugo.BackendUserID, cfg.Centrifugo.JwtToken)
	if err != nil {
		logrus.Fatal(ctx, err)
	}
	defer func() {
		closeErr := centrifugoClient.Close()
		logrus.Error(ctx, closeErr)
	}()

	twitchClient, err := twitch.GetClient(cfg.Twitch.Nick, cfg.Twitch.Pass)
	if err != nil {
		logrus.Fatal(ctx, err)
	}
	defer func() {
		closeErr := twitchClient.Close()
		logrus.Error(ctx, closeErr)
	}()

	allGames := make(map[string]*game.Game, len(registeredChannels))
	for _, ch := range registeredChannels {
		channelGame := game.InitGame()
		allGames[ch.Name] = &channelGame

		go game.Run(func(ctx context.Context, action game.Action) bool {
			err := centrifugoClient.Publish(cfg.Centrifugo.TwitchBossChannel, centrifugo.FormMessage(action.String()))
			if err != nil {
				logrus.Error(ctx, err)
				return false
			}
			return true
		}, &channelGame)
		go runTwitchListener(twitchClient, ch, channelGame.Channel)
	}

	webCfg := globalConfig.GetCfg().Web
	if err := fasthttp.ListenAndServe(webCfg.Host+":"+strconv.FormatInt(webCfg.Port, 10), fasthttp.CompressHandler(indexHandler)); err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func init() {
	var err error
	var environment string

	flag.StringVar(&environment, "env", "", "(re)generate config code")
	flag.Parse()

	env = config.GetEnvFromString(environment)
	cfg, err := config.Read(env, config.GetDefaultValuesPath())
	if err != nil {
		logrus.Fatal(err)
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = logrus.WarnLevel
		logrus.Warn(err)
	}

	logrus.SetLevel(logLevel)
	logrus.Debug("startup config", cfg)
	globalConfig.SetNew(cfg)

	content, err := ioutil.ReadFile(centrifugo.JSClientPath)
	if nil != err {
		logrus.Fatal(err)
	}
	jsClient = string(content)

	startTime = time.Now()

	go config.LiveRead(
		env,
		&globalConfig,
		func(e error) {
			logrus.Error(e)
		})
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	cfg := globalConfig.GetCfg()
	ctx.SetContentType("text/html; charset=utf8")
	ctx.Response.AppendBodyString(`
		<html>
		<script>
		` + jsClient + `
		</script>
		<script>
			var centrifuge = new Centrifuge('` + cfg.Centrifugo.URL + `');
			centrifuge.setToken("` + centrifugo.GetConnToken("112", cfg.Centrifugo.JwtToken, 0) + `");

			centrifuge.subscribe("` + cfg.Centrifugo.TwitchBossChannel + `", function(message) {
				console.log(message);
				document.getElementById('message-box').innerHTML += '<br>' + message.data.message;
			});

			centrifuge.connect();
		</script>
		<div id="message-box"></div>
		</html>
	`)
}

// Blocking. Subscribes to twitch chat, publishes messages to game channel
func runTwitchListener(
	twitchClient *twitch.Client,
	ch twitch.Channel,
	gameChannel chan<- game.Action,
) {
	ctx := context.Background()
	twitchClient.Listen(
		ch.Name,
		func(from, message string) { // message callback
			if action, err := ch.GetGameActionByViewer(from, message); err == nil {
				gameChannel <- action
			}
			logrus.Debug(from, ": ", message)
		},
		func(err error) { // error callback
			logrus.Fatal(ctx, err)
		},
		func(err error) { // warn callback
			logrus.Error(ctx, err)
		})
}
