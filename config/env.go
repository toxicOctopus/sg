package config

// Env type with hardcoded values because logic may be depended from env
type Env int

const (
	// Local development env
	Local Env = iota
	// Production environment
	Prod
)

var envMap map[Env]string

func init() {
	envMap = map[Env]string{
		Local: "local",
		Prod:  "prod",
	}
}

// GetEnvFromString get Env type
func GetEnvFromString(s string) Env {
	for env, str := range envMap {
		if str == s {
			return env
		}
	}

	return Local
}

// String Env to string
func (e Env) String() string {
	return envMap[e]
}
