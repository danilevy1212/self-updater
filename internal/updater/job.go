package updater

func (u *Updater) Run() {
	u.Logger.Info().
		Str("version", u.Meta.Version).
		Str("commit", u.Meta.Commit).
		Str("digest", u.Meta.Digest).
		Msg("Running updater job")
}
