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
        digest: $linux_amd64_digest,
        signatureBase64: $linux_amd64_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-linux-amd64")
      },
      {
        os: "linux",
        arch: "arm64",
        filename: "api-linux-arm64",
        digest: $linux_arm64_digest,
        signatureBase64: $linux_arm64_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-linux-arm64")
      },
      {
        os: "windows",
        arch: "amd64",
        filename: "api-windows-amd64.exe",
        digest: $windows_amd64_exe_digest,
        signatureBase64: $windows_amd64_exe_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-windows-amd64.exe")
      },
      {
        os: "windows",
        arch: "arm64",
        filename: "api-windows-arm64.exe",
        digest: $windows_arm64_exe_digest,
        signatureBase64: $windows_arm64_exe_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-windows-arm64.exe")
      },
      {
        os: "darwin",
        arch: "amd64",
        filename: "api-darwin-amd64",
        digest: $darwin_amd64_digest,
        signatureBase64: $darwin_amd64_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-darwin-amd64")
      },
      {
        os: "darwin",
        arch: "arm64",
        filename: "api-darwin-arm64",
        digest: $darwin_arm64_digest,
        signatureBase64: $darwin_arm64_sig,
        url: ($archiver_base_url + "/" + $archiver_owner + "/" + $archiver_repo + "/releases/download/" + $version + "/api-darwin-arm64")
      }
    ]
  }] + $old.versions)
}
