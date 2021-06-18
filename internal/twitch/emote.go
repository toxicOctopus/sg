package twitch

import "github.com/toxicOctopus/sg/internal/pkg/game"

type Emote struct {
	ID           int
	Name         string
	ImagePath    string      `json:"image_path"`
	ActionType   game.Action `json:"action_type"`
	ActionName   string      `json:"action_name"`
	ActionSource int         `json:"action_source"`
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
