package twitch

import "github.com/toxicOctopus/sg/game"

type Emote struct {
	Name       string
	ImagePath  string
	GameAction game.ViewerAction
}

type EmoteList []Emote

// MessageIsEmote checks if message is a valid emote from provided list
func MessageIsEmote(message string, list EmoteList) bool {
	for _, emote := range list {
		if emote.Name == message {
			return true
		}
	}

	return false
}
