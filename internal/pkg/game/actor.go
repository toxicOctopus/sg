package game

import "time"

type Viewer struct {
	ID             int
	Name           string
	HP             int
	BlockTimeStamp time.Time
	DodgeTimeStamp time.Time
}

type Host struct {
	HP             int
	BlockTimeStamp time.Time
	DodgeTimeStamp time.Time
}

// game variables which are dependant on viewer count
type variables struct {
	HostBlockTimeSpan time.Duration
	HostDodgeTimeSpan time.Duration

	ViewerBlockTimeSpan time.Duration
	ViewerDodgeTimeSpan time.Duration

	SweepDamage   int
	ExecuteDamage int

	ViewerAttackDamage    int
	ViewerOverPowerDamage int

	BossHP   int
	ViewerHP int
}

func (g *Game) calculateGameVariables() {
	g.Variables = variables{ //TODO implement game logic
		HostBlockTimeSpan:     time.Second,
		HostDodgeTimeSpan:     time.Second,
		ViewerBlockTimeSpan:   time.Second,
		ViewerDodgeTimeSpan:   time.Millisecond * 200,
		SweepDamage:           1,
		ExecuteDamage:         999,
		ViewerAttackDamage:    1,
		ViewerOverPowerDamage: 2,
		BossHP:                200,
		ViewerHP:              2,
	}
	g.Host.HP = g.Variables.BossHP
	for k, viewer := range g.Viewers {
		viewer.HP = g.Variables.ViewerHP
		g.Viewers[k] = viewer
	}
}

func (h *Host) IsAlive() bool {
	return h.HP > 0
}

func (v *Viewer) IsAlive() bool {
	return v.HP > 0
}