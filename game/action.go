package game

type Action int

type ViewerAction Action

type HostAction Action

const (
	ViewerDamage ViewerAction = iota + 1
	ViewerDodge
	ViewerBlock
	ViewerOP
)

const (
	HostSweep HostAction = iota + 300
	HostBlock
	HostTargetExecute
	HostFake
)
