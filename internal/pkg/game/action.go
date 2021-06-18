package game

type Action int

type ActionSource int

// Constants in the following block MUST be in the same order as "action_types" pg table. This is required for mapping.
const (
	ViewerDamage Action = iota + 1
	ViewerDodge
	ViewerBlock
	ViewerOverPower

	HostSweep
	HostBlock
	HostTargetExecute
	HostFake
)

const (
	HostAction ActionSource = iota +1
	ViewerAction
)
