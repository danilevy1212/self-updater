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
        digest: $lin_amd_digest,
        signatureBase64: $lin_amd_sig
      },
      {
        os: "linux",
        arch: "arm64",
        filename: "api-linux-arm64",
        digest: $lin_arm_digest,
        signatureBase64: $lin_arm_sig
      },
      {
        os: "windows",
        arch: "amd64",
        filename: "api-windows-amd64.exe",
        digest: $win_amd_digest,
        signatureBase64: $win_amd_sig
      },
      {
        os: "windows",
        arch: "arm64",
        filename: "api-windows-arm64.exe",
        digest: $win_arm_digest,
        signatureBase64: $win_arm_sig
      },
      {
        os: "darwin",
        arch: "amd64",
        filename: "api-darwin-amd64",
        digest: $mac_amd_digest,
        signatureBase64: $mac_amd_sig
      },
      {
        os: "darwin",
        arch: "arm64",
        filename: "api-darwin-arm64",
        digest: $mac_arm_digest,
        signatureBase64: $mac_arm_sig
      }
    ]
  }] + $old.versions)
}
