package models

import (
	"runtime"

	"github.com/danilevy1212/self-updater/internal/assets"
)

type ApplicationMeta struct {
	Digest                      string
	Version                     string
	Commit                      string
	IntegrityAuthorityPublicKey []byte
	OS                          string
	Arch                        string
}

func NewApplicationMeta(digest, version, commit string) ApplicationMeta {
	return ApplicationMeta{
		Digest:                      digest,
		Version:                     version,
		Commit:                      commit,
		IntegrityAuthorityPublicKey: assets.PublicKeyPEM,
		OS:                          runtime.GOOS,   // "linux", "windows" or "darwin"
		Arch:                        runtime.GOARCH, // "amd64" or "arm64"
	}
}
