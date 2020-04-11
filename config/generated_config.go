package config

type Config struct {
	ConfigReadInterval string `json:"configReadInterval"`
	LogLevel           string `json:"logLevel"`
	Web                struct {
		Host string `json:"host"`
		Port int64  `json:"port"`
	} `json:"web"`
	Ws struct {
		JwtToken string `json:"jwtToken"`
	} `json:"ws"`
}
