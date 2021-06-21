package twitch

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/toxicOctopus/sg/internal/pkg/game"
)

type Channel struct {
	ID       int
	Name     string
	Emotes   EmoteList // emotes usable for game
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

func (c *Channel) GetGameActionByViewer(viewerName, message string) (game.Action, error) {
	action := game.Action{
		Type:   0,
		Source: game.ViewerAction,
	}
	for _, emote := range c.Emotes {
		if emote.Name == message {
			action.Type = emote.ActionType
			action.ViewerName = viewerName
			return action, nil
		}
	}

	return action, errors.New("message doesn't fit the channel")
}
