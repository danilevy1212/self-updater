package models

type ReleaseManifest struct {
	Latest    string        `json:"latest"`
	PublicKey string        `json:"publicKey"`
	Versions  []ReleaseInfo `json:"versions"`
}

type ReleaseInfo struct {
	Version   string     `json:"version"`
	Commit    string     `json:"commit"`
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	Filename        string `json:"filename"`
	Digest          string `json:"digest"`
	SignatureBase64 string `json:"signatureBase64"`
	URL             string `json:"url"`
}
