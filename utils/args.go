package utils

type Env int

const (
	Dev Env = iota
	Test
	Prod
)

type Arguments struct {
	Host string `short:"h" long:"host" description:"host" default:"localhost"`
	Port string `short:"p" long:"port" description:"http port" default:"8080" optional:"y"`
	Env  string `short:"e" long:"env" description:"environment" default:"dev" optional:"y" choice:"dev" choice:"test" choice:"prod"`
}
