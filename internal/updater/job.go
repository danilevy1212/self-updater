package updater

func (u *Updater) Run() {
	u.Logger.Info().
		Str("version", u.Meta.Version).
		Str("commit", u.Meta.Commit).
		Str("digest", u.Meta.DigestString()).
		Str("public_key", string(u.Meta.IntegrityAuthorityPublicKey)).
		Msg("Running updater job")
}
