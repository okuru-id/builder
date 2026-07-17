// Package resources embeds static storefront assets (welcome page, etc.)
// so they ship inside the binary with no runtime file IO.
package resources

import _ "embed"

//go:embed welcome.html
var WelcomePage string
