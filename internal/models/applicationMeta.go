package models

import (
	"encoding/hex"
	"runtime"

	"github.com/danilevy1212/self-updater/internal/assets"
)

// NOTE  It is assumed the repo is public, in a closed sourced application, the fetching model would be different and this data may be unnecessary.
const (
	SourceName  = "self-updater"
	SourceOwner = "danilevy1212"
	Host        = "github.com"
)

type SourceInfo struct {
	Name  string
	Owner string
	Host  string
}

type ApplicationMeta struct {
	ExecutablePath   string
	Digest           []byte
	Version          string
	Commit           string
	AuthorsPublicKey []byte
	OS               string
	Arch             string
	SourceInfo       SourceInfo
}

func NewApplicationMeta(digest []byte, version, commit, exePath string) ApplicationMeta {
	return ApplicationMeta{
		Digest:           digest,
		Version:          version,
		Commit:           commit,
		AuthorsPublicKey: assets.PublicKeyPEM,
		OS:               runtime.GOOS,   // "linux", "windows" or "darwin"
		Arch:             runtime.GOARCH, // "amd64" or "arm64"
		ExecutablePath:   exePath,
		SourceInfo: SourceInfo{
			Name:  SourceName,
			Owner: SourceOwner,
			Host:  Host,
		},
	}
}

func (am ApplicationMeta) DigestString() string {
	return hex.EncodeToString(am.Digest)
}
