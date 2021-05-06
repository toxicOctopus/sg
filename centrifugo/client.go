package centrifugo

import (
	"github.com/centrifugal/centrifuge-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	JSClientPath = "resources/wsClient.js"
)

type eventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	logrus.Debug("centrifugo: connected")
}

func (h *eventHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	logrus.Error("centrifugo: error ", e.Message)
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	logrus.Error("centrifugo: disconnected ", e.Reason)
}

func GetClient(wsURL, userID, token string) (*centrifuge.Client, error) {
	c := centrifuge.New(wsURL, centrifuge.DefaultConfig())
	c.SetToken(GetConnToken(userID, token, 0))
	handler := &eventHandler{}
	c.OnDisconnect(handler)
	c.OnConnect(handler)
	c.OnError(handler)

	err := c.Connect()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetConnToken(user, token string, exp int64) string {
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(token))
	if err != nil {
		panic(err)
	}
	return t
}

// forms message usable in app for centrifugo
func FormMessage(message string) []byte {
	return []byte("{\"message\":\"" + message + "\"}")
}
