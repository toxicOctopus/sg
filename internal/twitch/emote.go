package twitch

import "github.com/toxicOctopus/sg/internal/pkg/game"

type Emote struct {
	Name       string
	ImagePath  string
	GameAction game.ViewerAction
}

type EmoteList []Emote

// MessageIsEmote checks if message is a valid emote from provided list
func (l EmoteList) MessageIsEmote(message string) bool {
	for _, emote := range l {
		if emote.Name == message {
			return true
		}
	}

	return false
}