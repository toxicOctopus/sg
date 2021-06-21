package game

import (
	"time"

	"github.com/toxicOctopus/sg/pkg/timer"
)

type Cooldown int

const (
	Sweep Cooldown = iota + 1
	Block
	Execute
)

const (
	SweepCD   = time.Second * 5
	BlockCD   = time.Second * 2
	ExecuteCD = time.Second * 15
)

type Cooldowns struct {
	SweepTimer   *timer.Timer
	BlockTimer   *timer.Timer
	ExecuteTimer *timer.Timer
}
