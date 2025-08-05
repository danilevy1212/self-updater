package models

type ApplicationMeta struct {
	Digest                      string
	Version                     string
	Commit                      string
	IntegrityAuthorityPublicKey []byte
}
