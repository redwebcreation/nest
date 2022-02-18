package service

import "sort"

type EnvMap map[string]string

func (e EnvMap) ForDocker() []string {
	env := make([]string, 0, len(e))
	for k, v := range e {
		env = append(env, k+"="+v)
	}

	// sort the env to make it deterministic
	sort.Strings(env)

	return env
}
