// Package integration provides embedded shell integration snippets.
package integration

import _ "embed"

// Bash contains the bash shell integration script.
//
//go:embed bash.sh
var Bash string

// Zsh contains the zsh shell integration script.
//
//go:embed zsh.sh
var Zsh string

// ZshFzf contains the zsh shell integration script with fzf support.
//
//go:embed zsh-fzf.sh
var ZshFzf string
