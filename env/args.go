package env

type Env int

const (
	Dev Env = iota
	Test
	Prod
)

type ConnectionString struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type Config struct {
	Web struct {
		ConnectionString
	} `json:"web"`
}
