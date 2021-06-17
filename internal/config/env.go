package config

type Env int

const (
	Local Env = iota
	Test
	Production
)

var envMap map[Env]string

func init() {
	envMap = map[Env]string{
		Local:      "local",
		Test:       "test",
		Production: "production",
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
