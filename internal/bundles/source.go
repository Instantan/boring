package bundles

import _ "embed"

//go:embed scss/out.es5.js
var scssSource []byte

func SCSSSource() []byte {
	return scssSource
}
