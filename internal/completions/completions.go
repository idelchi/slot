// Package completions provides embedded shell completion snippets.
package completions

import _ "embed"

// Bash contains the bash shell completion script.
//
//go:embed bash.sh
var Bash string

// Zsh contains the zsh shell completion script.
//
//go:embed zsh.sh
var Zsh string
