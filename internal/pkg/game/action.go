package game

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/toxicOctopus/sg/pkg/common"
)

// Constants in the following block MUST be in the same order as "action_types" pg table. This is required for mapping.
// TODO test for that
const (
	ViewerDamage ActionType = iota + 1
	ViewerDodge
	ViewerBlock
	ViewerOverPower

	HostSweep
	HostBlock
	HostTargetExecute
	HostFake

	RegisterViewer

	HostGameStart
	HostGamePause
	HostGameStop

	GameStop

	CooldownRefreshed

	ViewerDamaged
	ViewerDead

	HostBlocked
	HostDead
)

const (
	HostAction ActionSource = iota + 1
	ViewerAction
	GameAction
)

type ActionType int

type ActionSource int

type Action struct {
	DamageDealt    int
	Type           ActionType
	Source         ActionSource
	ViewerName     string
	VictimNickname string
	Cooldown       Cooldown
}

var (
	affectingGameTime = []int{
		int(HostGameStart),
		int(HostGamePause),
		int(HostGameStop),
		int(GameStop),
	}
)

func (a *ActionType) IsAffectingGameTime() bool {
	if common.IntInSlice(int(*a), affectingGameTime) {
		return true
	}

	return false
}

func (ca Action) String() string {
	str, err := json.Marshal(ca)
	if err != nil {
		logrus.Error("building complex action", ca.Type, err)
	}

	return string(str)
}
