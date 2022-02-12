package cloud

func AddToken(token string) error {
	Config.Tokens = append(Config.Tokens, token)

	return Config.Save()
}
