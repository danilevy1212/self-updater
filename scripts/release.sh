#!/usr/bin/env bash
set -euo pipefail

APP_NAME="api"
BIN_DIR="bin"
MANIFEST="internal/assets/release.json"
SIGN_CMD="go run ./cmd/sign"

VERSION="${VERSION:-unknown}"
COMMIT="${COMMIT:-unknown}"
SIGNING_KEY_ENV="${SIGNING_KEY_ENV:-SIGNING_KEY_PEM}"
PUBLIC_KEY="${PUBLIC_KEY:-}"

SIGN_KEY_FILE="$(mktemp)"
trap 'rm -f "$SIGN_KEY_FILE"' EXIT

if [[ -z "${!SIGNING_KEY_ENV:-}" ]]; then
  echo "FATAL: Env variable \$${SIGNING_KEY_ENV} is not set."
  exit 1
fi

echo "${!SIGNING_KEY_ENV}" > "$SIGN_KEY_FILE"

targets=(
  "linux-amd64"
  "linux-arm64"
  "windows-amd64.exe"
  "windows-arm64.exe"
  "darwin-amd64"
  "darwin-arm64"
)

declare -A DIGESTS
declare -A SIGS

for target in "${targets[@]}"; do
  bin="$BIN_DIR/$APP_NAME-$target"
  digest=$(sha256sum "$bin" | cut -d ' ' -f1)
  sig=$($SIGN_CMD "$SIGN_KEY_FILE" "$bin")
  key="${target//[^a-zA-Z0-9]/_}"
  DIGESTS[$key]="$digest"
  SIGS[$key]="$sig"
done

TMP_MANIFEST=$(mktemp)
jq_args=(
  --arg version "$VERSION"
  --arg commit "$COMMIT"
  --arg pubkey "$PUBLIC_KEY"
)

for target in "${!DIGESTS[@]}"; do
  digest=${DIGESTS[$target]}
  sig=${SIGS[$target]}
  jq_args+=(--arg "${target}_digest" "$digest")
  jq_args+=(--arg "${target}_sig" "$sig")
done

jq "${jq_args[@]}" -f scripts/merge_manifest.jq "$MANIFEST" > "$TMP_MANIFEST"
mv "$TMP_MANIFEST" "$MANIFEST"

$SIGN_CMD "$SIGN_KEY_FILE" "$MANIFEST" > "$MANIFEST.sig.base64"

echo "Manifest and signature updated: $MANIFEST"
