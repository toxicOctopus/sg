package game

import (
	"context"
	"time"

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

	Host      Host
	Cooldowns Cooldowns

	Channel   chan Action
	Viewers   map[string]Viewer
	Variables variables
}

func InitGame() Game {
	channel := make(chan Action)
	return Game{
		State: AwaitingStart,

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
	sendAction func(ctx context.Context, action Action) bool,
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
				if !game.IsRunnable() {
					continue
				}
				viewer, exist := game.Viewers[msg.ViewerName]
				if !exist || !viewer.IsAlive() {
					continue
				}
				now := time.Now()
				hostBlocked := now.Sub(game.Host.BlockTimeStamp) <= game.Variables.HostBlockTimeSpan
				if hostBlocked {
					sendAction(ctx, Action{
						Type:        HostBlocked,
						Source:      GameAction,
						DamageDealt: 0,
					})
					continue
				}

				sendAction(ctx, Action{
					Type:        ViewerDamage,
					Source:      ViewerAction,
					ViewerName:  msg.ViewerName,
					DamageDealt: game.Variables.ViewerAttackDamage,
				})
				game.Host.HP = game.Host.HP - game.Variables.ViewerAttackDamage
				if !game.Host.IsAlive() {
					sendAction(ctx, Action{
						Type:   HostDead,
						Source: GameAction,
					})
					game.State = Ended
					newGame := InitGame()
					game = &newGame
					Run(sendAction, game)
				}
			case ViewerBlock:
				if !game.IsRunnable() {
					continue
				}
				viewer, exist := game.Viewers[msg.ViewerName]
				if !exist || !viewer.IsAlive() {
					continue
				}
				viewer.BlockTimeStamp = time.Now()
				game.Viewers[msg.ViewerName] = viewer

				sendAction(ctx, Action{
					Type:       ViewerBlock,
					Source:     ViewerAction,
					ViewerName: msg.ViewerName,
				})

			case ViewerDodge:
				if !game.IsRunnable() {
					continue
				}
				viewer, exist := game.Viewers[msg.ViewerName]
				if !exist || !viewer.IsAlive() {
					continue
				}
				viewer.DodgeTimeStamp = time.Now()
				game.Viewers[msg.ViewerName] = viewer

				sendAction(ctx, Action{
					Type:       ViewerDodge,
					Source:     ViewerAction,
					ViewerName: msg.ViewerName,
				})

			case ViewerOverPower:
				if !game.IsRunnable() {
					continue
				}
				viewer, exist := game.Viewers[msg.ViewerName]
				if !exist || !viewer.IsAlive() {
					continue
				}
				//TODO implementation ????
			case HostSweep:
				if !game.IsRunnable() || game.Cooldowns.SweepTimer.IsActive() {
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
				sendAction(ctx, Action{
					Type:   HostSweep,
					Source: HostAction,
				})
				game.Cooldowns.SweepTimer.Start()

				for name, viewer := range game.Viewers {
					viewerDodged := now.Sub(viewer.DodgeTimeStamp) <= game.Variables.ViewerDodgeTimeSpan
					viewerBlocked := now.Sub(viewer.BlockTimeStamp) <= game.Variables.ViewerBlockTimeSpan
					if !viewerDodged {
						damageDealt := 0
						if viewerBlocked {
							damageDealt = game.Variables.SweepDamage / 2
						} else {
							damageDealt = game.Variables.SweepDamage
						}
						viewer.HP = viewer.HP - damageDealt
						game.Viewers[name] = viewer

						sendAction(ctx, Action{
							Type:           ViewerDamaged,
							VictimNickname: viewer.Name,
							Source:         GameAction,
							DamageDealt:    damageDealt,
						})
					}

					if !viewer.IsAlive() {
						delete(game.Viewers, name)
						sendAction(ctx, Action{
							Type:           ViewerDead,
							VictimNickname: viewer.Name,
							Source:         GameAction,
						})
					}
				}
			case HostBlock:
				if !game.IsRunnable() || game.Cooldowns.BlockTimer.IsActive() {
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
				game.Cooldowns.BlockTimer.Start()
				sendAction(ctx, Action{
					Type:   HostBlock,
					Source: HostAction,
				})
			case HostTargetExecute:
				if !game.IsRunnable() || game.Cooldowns.ExecuteTimer.IsActive() || len(game.Viewers) == 1 {
					continue
				}

				viewer, exist := game.Viewers[msg.VictimNickname]
				if !exist || !viewer.IsAlive() {
					continue
				}

				game.Cooldowns.ExecuteTimer = timer.AfterFunc(ExecuteCD, func() {
					game.Channel <- Action{
						Type:     CooldownRefreshed,
						Source:   GameAction,
						Cooldown: Execute,
					}
				})
				game.Cooldowns.ExecuteTimer.Start()
				sendAction(ctx, Action{
					Type:           HostTargetExecute,
					Source:         HostAction,
					VictimNickname: msg.VictimNickname,
				})

				viewer.HP = viewer.HP - game.Variables.ExecuteDamage
				game.Viewers[msg.VictimNickname] = viewer

				sendAction(ctx, Action{
					Type:           ViewerDamaged,
					VictimNickname: viewer.Name,
					Source:         GameAction,
					DamageDealt:    game.Variables.ExecuteDamage,
				})

				if !viewer.IsAlive() {
					delete(game.Viewers, viewer.Name)
					sendAction(ctx, Action{
						Type:           ViewerDead,
						VictimNickname: viewer.Name,
						Source:         GameAction,
					})
				}

			case RegisterViewer:
				if _, exist := game.Viewers[msg.ViewerName]; !exist {
					game.Viewers[msg.ViewerName] = Viewer{
						Name: msg.ViewerName,
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
					game.calculateGameVariables()
				}
				game.State = Running
			case HostGameStop:
				fallthrough
			case GameStop:
				game.State = Ended
				newGame := InitGame()
				game = &newGame
				Run(sendAction, game)

			case CooldownRefreshed:
				sendAction(ctx, Action{
					Type:     CooldownRefreshed,
					Source:   GameAction,
					Cooldown: msg.Cooldown,
				})
			}
		}
	}
}

func (g *Game) IsRunnable() bool {
	if g.State != Running && g.Host.IsAlive() {
		return false
	}

	return true
}
