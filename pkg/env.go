package pkg

type EnvMap map[string]string

func (e EnvMap) ToDockerEnv() []string {
	env := make([]string, 0, len(e))
	for k, v := range e {
		env = append(env, k+"="+v)
	}
	return env
}
