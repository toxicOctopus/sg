package utils

type Env int

const (
	Dev Env = iota
	Test
	Prod
)

type Arguments struct {
	Env string `short:"e" long:"env" description:"environment" default:"dev" optional:"y" choice:"dev" choice:"test" choice:"prod"`
}