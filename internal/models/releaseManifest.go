package models

import "errors"

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

func (rm *ReleaseManifest) GetVersionInfo(version string) (*ReleaseInfo, error) {
	for _, v := range rm.Versions {
		if v.Version == version {
			return &v, nil
		}
	}

	return nil, errors.New("version not found in manifest")
}

func (ri *ReleaseInfo) GetArtifactForPlatform(os, arch string) (*Artifact, error) {
	for _, a := range ri.Artifacts {
		if a.OS == os && a.Arch == arch {
			return &a, nil
		}
	}

	return nil, errors.New("artifact not found for the specified platform")
}
