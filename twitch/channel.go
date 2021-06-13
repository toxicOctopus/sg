package twitch

import (
	"github.com/pkg/errors"
	"github.com/toxicOctopus/sg/game"
	"strings"
	"time"
)

type Channel struct {
	Name     string
	Emotes   EmoteList // emotes usable for game
	ActionCD time.Duration
}

type RegisteredChannels struct {
	Channels []Channel
}

func (rc RegisteredChannels) GetChannel(name string) (Channel, error) {
	for _, channel := range rc.Channels {
		if strings.ToLower(channel.Name) == strings.ToLower(name) {
			return channel, nil
		}
	}

	return Channel{}, errors.New("channel not found")
}

//TODO implement
func LoadRegisteredChannels() RegisteredChannels {
	return RegisteredChannels{
		Channels: []Channel{
			{
				Name: "toxic_octopuz",
				Emotes: EmoteList{
					{
						Name:       "Kappa",
						ImagePath:  "",
						GameAction: game.ViewerDodge,
					},
					{
						Name:       "SMOrc",
						ImagePath:  "",
						GameAction: game.ViewerDamage,
					},
					{
						Name:       "4Head",
						ImagePath:  "",
						GameAction: game.ViewerBlock,
					},
				},
				ActionCD: time.Second,
			},
		},
	}
}
