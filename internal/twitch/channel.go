package twitch

import (
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Channel struct {
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