package twitch

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Channel struct {
	ID       int
	Name     string
	Emotes   EmoteList // emotes usable for game
	ActionCD time.Duration
}

type RegisteredChannels []Channel

func (rc RegisteredChannels) GetChannel(name string) (Channel, error) {
	for _, channel := range rc {
		if strings.ToLower(channel.Name) == strings.ToLower(name) {
			return channel, nil
		}
	}

	return Channel{}, errors.New("channel not found")
}
