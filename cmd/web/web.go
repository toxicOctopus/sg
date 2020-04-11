package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/toxicOctopus/sg/config"
	"github.com/valyala/fasthttp"
)

const (
	twitchBossChannel = "public:tb"
)

var (
	jsClient     string
	startTime    time.Time
	env          config.Env
	globalConfig config.LiveConfig
)

func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(globalConfig.GetCfg().Ws.JwtToken))
	if err != nil {
		panic(err)
	}
	return t
}

type eventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	log.Println("Connected")
}

func (h *eventHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	log.Println("Error", e.Message)
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Println("Disconnected", e.Reason)
}

func newConnection() *centrifuge.Client {
	wsURL := "ws://localhost:8000/connection/websocket"

	c := centrifuge.New(wsURL, centrifuge.DefaultConfig())
	c.SetToken(connToken("551", 0))
	handler := &eventHandler{}
	c.OnDisconnect(handler)
	c.OnConnect(handler)
	c.OnError(handler)

	err := c.Connect()
	if err != nil {
		logrus.Fatalln(err)
	}
	return c
}

func main() {
	logrus.Info(startTime)
	logrus.Info("Start program")
	c := newConnection()
	defer func() {
		closeErr := c.Close()
		logrus.Error(closeErr)
	}()

	err := c.Publish(twitchBossChannel, []byte("{\"kek\":\"lul\"}"))
	if err != nil {
		logrus.Fatal(err)
	}

	webCfg := globalConfig.GetCfg().Web
	if err := fasthttp.ListenAndServe(webCfg.Host + ":" + strconv.FormatInt(webCfg.Port, 10), fasthttp.CompressHandler(indexHandler)); err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func init() {
	var err error

	var environment string
	flag.StringVar(&environment, "env", "", "(re)generate config code")
	flag.Parse()

	env = config.GetEnvFromString(environment)
	cfg, err := config.Read(env)
	if err != nil {
		logrus.Fatal(err)
	}
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = logrus.WarnLevel
		logrus.Warn(err)
	}
	logrus.SetLevel(logLevel)
	logrus.Debug("result config", cfg)
	globalConfig.SetNew(cfg)

	content, err := ioutil.ReadFile("resources/wsClient.js")
	if nil != err {
		logrus.Fatal(err)
	}
	jsClient = string(content)
	startTime = time.Now()

	go config.LiveRead(env, &globalConfig, config.StringToUpdateInterval(globalConfig.GetCfg().ConfigReadInterval))
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf8")
	ctx.Response.AppendBodyString(`
		<html>
		<script>
		` + jsClient + `
		</script>
		<script>
			var centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket');
			centrifuge.setToken("` + connToken("112", 0) + `");

			centrifuge.subscribe("` + twitchBossChannel + `", function(message) {
				console.log(message);
			});

			centrifuge.connect();
		</script>
		ty pidor
		</html>
	`)
}
