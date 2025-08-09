package fixtures

import _ "embed"

//go:embed release_fixture_json_sig_base64.txt
var ReleaseSignatureFixture []byte
