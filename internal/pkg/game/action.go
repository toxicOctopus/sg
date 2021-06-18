package game

type Action int

type ActionSource int

type ComplexAction struct {
	Action         Action
	Source         ActionSource
	VictimNickname string
	Cooldown       Cooldown
}

// Constants in the following block MUST be in the same order as "action_types" pg table. This is required for mapping.
// TODO test for that
const (
	ViewerDamage Action = iota + 1
	ViewerDodge
	ViewerBlock
	ViewerOverPower

	HostSweep
	HostBlock
	HostTargetExecute

	HostGameStart
	HostGamePause
	HostGameStop

	GameStop
	CooldownRefreshed
)

const (
	HostAction ActionSource = iota + 1
	ViewerAction
)
