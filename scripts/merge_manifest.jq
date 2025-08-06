. as $old | {
  latest: $version,
  publicKey: ($old.publicKey // $pubkey),
  versions: ([{
    version: $version,
    commit: $commit,
    artifacts: [{
      os: "linux",
      arch: "amd64",
      filename: $filename,
      digest: $digest,
      signatureBase64: $sig
    }]
  }] + $old.versions)
}
