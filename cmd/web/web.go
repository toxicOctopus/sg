package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/jessevdk/go-flags"
	"github.com/valyala/fasthttp"

	"sg/env"
)

var jsClient string

const (
	TwitchBossChannel = "public:tb"
)

type arguments struct {
	env.Arguments
}

var (
	args      arguments
	startTime time.Time
	
	JWTToken string
)


func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWTToken))
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
		log.Fatalln(err)
	}
	return c
}


func main() {
	log.Println(startTime)
	log.Println("Start program")
	c := newConnection()
	defer c.Close()


	c.Publish(TwitchBossChannel, []byte("{\"kek\":\"lul\"}"))

	if err := fasthttp.ListenAndServe(args.Host + ":" + args.Port, fasthttp.CompressHandler(indexHandler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func init() {
	_, err := flags.Parse(&args)
	if nil != err {
		os.Exit(1)
	}
	content, err := ioutil.ReadFile("resources/wsClient.js")
	if nil != err {
		os.Exit(1)
	}
	jsClient = string(content)

	startTime = time.Now()
	JWTToken = "68a91e24-4a3f-4046-b8c7-faccc884f9fc"
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf8")
	ctx.Response.AppendBodyString(`
		<html>
		<script>
		`+jsClient+`
		</script>
		<script>
			var centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket');
			centrifuge.setToken("`+connToken("112", 0)+`");

			centrifuge.subscribe("`+TwitchBossChannel+`", function(message) {
				console.log(message);
			});

			centrifuge.connect();
		</script>
		ty pidor
		</html>
	`)
}