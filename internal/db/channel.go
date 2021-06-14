package db

import (
	"time"

	"github.com/toxicOctopus/sg/internal/pkg/game"
	"github.com/toxicOctopus/sg/internal/twitch"
)

//TODO implement
func LoadRegisteredChannels() twitch.RegisteredChannels {
	return twitch.RegisteredChannels{
		{
			Name: "toxic_octopuz",
			Emotes: twitch.EmoteList{
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
	}
}
