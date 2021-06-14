package config

type Config struct {
	Centrifugo struct {
		BackendUserID     string `json:"backendUserID"`
		JwtToken          string `json:"jwtToken"`
		TwitchBossChannel string `json:"twitchBossChannel"`
		URL               string `json:"url"`
	} `json:"centrifugo"`
	ConfigReadInterval string `json:"configReadInterval"`
	LogLevel           string `json:"logLevel"`
	Twitch             struct {
		Nick string `json:"nick"`
		Pass string `json:"pass"`
	} `json:"twitch"`
	Web struct {
		Host string `json:"host"`
		Port int64  `json:"port"`
	} `json:"web"`
}
