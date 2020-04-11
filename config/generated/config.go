package config

type Config struct {
	Web struct {
		Host string `json:"host"`
		Port int64  `json:"port"`
	} `json:"web"`
}
