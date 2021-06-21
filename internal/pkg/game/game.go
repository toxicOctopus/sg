package game

import (
	"context"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/sirupsen/logrus"

	"github.com/toxicOctopus/sg/internal/centrifugo"
	"github.com/toxicOctopus/sg/internal/twitch"
	"github.com/toxicOctopus/sg/pkg/timer"
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
	State State

	Rules twitch.EmoteList

	Host      Host
	Cooldowns Cooldowns

	Channel   chan Action
	Viewers   map[string]Viewer
	Variables variables
}

func InitGame(twitchChannel twitch.Channel) Game {
	channel := make(chan Action)
	return Game{
		State: AwaitingStart,

		Rules: twitchChannel.Emotes,

		Host: Host{
			HP:             100,
			BlockTimeStamp: time.Time{},
			DodgeTimeStamp: time.Time{},
		},

		Channel: channel,
	}
}

// blocking call
func Run(
	twitchClient *twitch.Client,
	centrifugoClient *centrifuge.Client,
	twitchChannel twitch.Channel,
	centrifugoChannel string,
	game *Game,
) {
	ctx := context.Background()
	for {
		select {
		case msg := <-game.Channel:
			if game.State == Pause && !msg.Type.IsAffectingGameTime() {
				continue
			}

			switch msg.Type {
			case ViewerDamage:
				if !game.Runnable() {
					continue
				}
				now := time.Now()
				blocked := now.Sub(game.Host.BlockTimeStamp) <= game.Variables.HostBlockTimeSpan
				if blocked {
					err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
						Type:        HostBlocked,
						Source:      GameAction,
						DamageDealt: 0,
					}.String()))
					if err != nil {
						logrus.Error(err)
					}
					continue
				}

				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:        ViewerDamage,
					Source:      ViewerAction,
					ViewerName:  msg.ViewerName,
					DamageDealt: game.Variables.ViewerAttackDamage,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				game.Host.HP = game.Host.HP - game.Variables.ViewerAttackDamage
				if game.Host.HP <= 0 {
					err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
						Type:   HostDead,
						Source: GameAction,
					}.String()))
					if err != nil {
						logrus.Error(err)
					}
					game.State = Ended
				}
			case ViewerBlock:
				if !game.Runnable() {
					continue
				}
				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:       ViewerBlock,
					Source:     ViewerAction,
					ViewerName: msg.ViewerName,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				if viewer, alive := game.Viewers[msg.ViewerName]; alive {
					viewer.BlockTimeStamp = time.Now()
					game.Viewers[msg.ViewerName] = viewer
				}
			case ViewerDodge:
				if !game.Runnable() {
					continue
				}
				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:       ViewerDodge,
					Source:     ViewerAction,
					ViewerName: msg.ViewerName,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				if viewer, alive := game.Viewers[msg.ViewerName]; alive {
					viewer.DodgeTimeStamp = time.Now()
					game.Viewers[msg.ViewerName] = viewer
				}
			case ViewerOverPower:
				if !game.Runnable() {
					continue
				}
				//TODO implementation ????
			case HostSweep:
				if !game.Runnable() || game.Cooldowns.SweepTimer.IsActive() {
					continue
				}
				now := time.Now()

				game.Cooldowns.SweepTimer = timer.AfterFunc(SweepCD, func() {
					game.Channel <- Action{
						Type:     CooldownRefreshed,
						Source:   GameAction,
						Cooldown: Sweep,
					}
				})
				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:   HostSweep,
					Source: HostAction,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				game.Cooldowns.SweepTimer.Start()

				for name, viewer := range game.Viewers {
					dodged := now.Sub(viewer.DodgeTimeStamp) <= game.Variables.ViewerDodgeTimeSpan
					blocked := now.Sub(viewer.BlockTimeStamp) <= game.Variables.ViewerBlockTimeSpan
					if !dodged {
						damageDealt := 0
						if blocked {
							damageDealt = game.Variables.SweepDamage / 2
						} else {
							damageDealt = game.Variables.SweepDamage
						}
						viewer.HP = viewer.HP - damageDealt
						game.Viewers[name] = viewer

						err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
							Type:           ViewerDamaged,
							VictimNickname: viewer.Name,
							Source:         GameAction,
							DamageDealt:    damageDealt,
						}.String()))
						if err != nil {
							logrus.Error(err)
						}
					}

					if viewer.HP <= 0 {
						delete(game.Viewers, name)
						err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
							Type:           ViewerDead,
							VictimNickname: viewer.Name,
							Source:         GameAction,
						}.String()))
						if err != nil {
							logrus.Error(err)
						}
					}
				}
			case HostBlock:
				if !game.Runnable() || game.Cooldowns.BlockTimer.IsActive() {
					continue
				}
				game.Cooldowns.BlockTimer = timer.AfterFunc(BlockCD, func() {
					game.Channel <- Action{
						Type:     CooldownRefreshed,
						Source:   GameAction,
						Cooldown: Block,
					}
				})

				game.Host.BlockTimeStamp = time.Now()

				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:   HostBlock,
					Source: HostAction,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				game.Cooldowns.BlockTimer.Start()
			case HostTargetExecute:
				if !game.Runnable() || game.Cooldowns.ExecuteTimer.IsActive() || len(game.Viewers) == 1 {
					continue
				}

				game.Cooldowns.ExecuteTimer = timer.AfterFunc(ExecuteCD, func() {
					game.Channel <- Action{
						Type:     CooldownRefreshed,
						Source:   GameAction,
						Cooldown: Execute,
					}
				})
				err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:           HostTargetExecute,
					Source:         HostAction,
					VictimNickname: msg.VictimNickname,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}
				game.Cooldowns.ExecuteTimer.Start()

				viewer, alive := game.Viewers[msg.VictimNickname]
				if !alive {
					continue
				}

				viewer.HP = viewer.HP - game.Variables.ExecuteDamage
				game.Viewers[msg.VictimNickname] = viewer

				err = centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
					Type:           ViewerDamaged,
					VictimNickname: viewer.Name,
					Source:         GameAction,
					DamageDealt:    game.Variables.ExecuteDamage,
				}.String()))
				if err != nil {
					logrus.Error(err)
				}

				if viewer.HP <= 0 {
					delete(game.Viewers, viewer.Name)
					err := centrifugoClient.Publish(centrifugoChannel, centrifugo.FormMessage(Action{
						Type:           ViewerDead,
						VictimNickname: viewer.Name,
						Source:         GameAction,
					}.String()))
					if err != nil {
						logrus.Error(err)
					}
				}

			case HostGamePause:
				game.State = Pause
				if game.Cooldowns.SweepTimer.IsActive() {
					game.Cooldowns.SweepTimer.Pause()
				}
				if game.Cooldowns.BlockTimer.IsActive() {
					game.Cooldowns.BlockTimer.Pause()
				}
				if game.Cooldowns.ExecuteTimer.IsActive() {
					game.Cooldowns.ExecuteTimer.Pause()
				}
			case HostGameStart:
				if game.State == Pause {
					if game.Cooldowns.SweepTimer.IsPaused() {
						game.Cooldowns.SweepTimer.Start()
					}
					if game.Cooldowns.BlockTimer.IsPaused() {
						game.Cooldowns.BlockTimer.Start()
					}
					if game.Cooldowns.ExecuteTimer.IsPaused() {
						game.Cooldowns.ExecuteTimer.Start()
					}
				} else {
					game.State = Starting
					viewers, err := twitchClient.GetViewers(twitchChannel.Name)
					if err != nil {
						logrus.Error(ctx, err)
						game.State = AwaitingStart
						continue
					}
					game.Viewers = viewers
					game.calculateGameVariables()
				}
				game.State = Running
			case HostGameStop:
				fallthrough
			case GameStop:
				game.State = Ended
				newGame := InitGame(twitchChannel)
				game = &newGame
				Run(twitchClient, centrifugoClient, twitchChannel, centrifugoChannel, game)

			case CooldownRefreshed:
				//TODO push to game centrifugo channel
			}
		}
	}
}

func (g *Game) Runnable() bool {
	if g.State != Running && g.Host.IsAlive() {
		return false
	}

	return true
}
