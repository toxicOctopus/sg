package game

import (
	"context"
	"sync"
	"time"

	"github.com/toxicOctopus/sg/internal/twitch"
)

type State int

const (
	AwaitingStart State = iota + 1
	Starting
	Running
	Pause
	Ended
)

type Game struct {
	State      State
	StateMutex *sync.Mutex

	Rules    twitch.EmoteList
	ActionCD time.Duration

	BossHP    int
	Cooldowns Cooldowns

	Channel chan ComplexAction
	Viewers map[string]Viewer
}

type Cooldown int64

type Cooldowns struct { // cooldown start timestamps
	Sweep   Cooldown
	Block   Cooldown
	Execute Cooldown
}

func InitGame(twitchChannel twitch.Channel) Game {
	channel := make(chan ComplexAction)
	return Game{
		State:      AwaitingStart,
		StateMutex: &sync.Mutex{},

		Rules:    twitchChannel.Emotes,
		ActionCD: twitchChannel.ActionCD,

		BossHP: 100,
		Cooldowns: Cooldowns{
			Sweep:   0,
			Block:   0,
			Execute: 0,
		},
		Channel: channel,
	}
}

// blocking code
func Run(ctx context.Context, twitchChannel twitch.Channel, game *Game) {
	for {
		select {
		case msg := <-game.Channel:
			switch msg.Action {
			case ViewerDamage:
			case ViewerBlock:
			case ViewerDodge:
			case ViewerOverPower:

			case HostSweep:
			case HostBlock:
			case HostTargetExecute:

			case HostGameStart:
				game.State = Starting
				//TODO readviewers
				game.State = Running
			case HostGamePause:
				game.State = Pause
			case HostGameStop:
				fallthrough
			case GameStop:
				game.State = Ended
				newGame := InitGame(twitchChannel)
				game = &newGame
				Run(ctx, twitchChannel, game)

			case CooldownRefreshed:
			}
		}
	}
}
