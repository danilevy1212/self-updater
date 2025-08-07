. as $old | {
  latest: $version,
  publicKey: ($old.publicKey // $pubkey),
  versions: ([{
    version: $version,
    commit: $commit,
    artifacts: [
      {
        os: "linux",
        arch: "amd64",
        filename: "api-linux-amd64",
        digest: $lin_digest,
        signatureBase64: $lin_sig
      },
      {
        os: "windows",
        arch: "amd64",
        filename: "api-windows-amd64.exe",
        digest: $win_digest,
        signatureBase64: $win_sig
      },
      {
        os: "darwin",
        arch: "amd64",
        filename: "api-darwin-amd64",
        digest: $mac_digest,
        signatureBase64: $mac_sig
      }
    ]
  }] + $old.versions)
}
