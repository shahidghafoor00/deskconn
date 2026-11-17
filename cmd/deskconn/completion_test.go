package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixBashCompletion(t *testing.T) {
	input := `_desk_bash_autocomplete() {
    local cur prev opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( ${COMP_WORDS[0]} --completion-bash "${COMP_WORDS[@]:1:$COMP_CWORD}" )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}
complete -F _desk_bash_autocomplete -o default desk
`

	out := fixBashCompletionCommands(fixBashCompletionSpacing(fixBashCompletionWordBreaks(input)))

	assert.Contains(t, out, "COMP_WORDBREAKS=${COMP_WORDBREAKS//:}")
	assert.Contains(t, out, "compopt -o filenames")
	assert.Contains(t, out, "compopt -o nospace")
	assert.Contains(t, out, "complete -F _desk_bash_autocomplete -o default desk deskconn")
}

func TestFixZshCompletion(t *testing.T) {
	input := `#compdef desk

_desk() {
    local matches=($(${words[1]} --completion-bash "${(@)words[2,$CURRENT]}"))
    compadd -a matches

    if [[ $compstate[nmatches] -eq 0 && $words[$CURRENT] != -* ]]; then
        _files
    fi
}

if [[ "$(basename -- ${(%):-%x})" != "_desk" ]]; then
    compdef _desk desk
fi
`

	out := fixZshCompletion(input)

	assert.Contains(t, out, "#compdef desk deskconn")
	assert.Contains(t, out, "compdef _desk desk deskconn")
	assert.Contains(t, out, `lines=("${(@f)$(${words[1]} --completion-bash "${(@)words[2,$CURRENT]}")}")`)
	assert.Contains(t, out, "compadd -S '' -a dirs_and_colons")
	assert.Contains(t, out, "compadd -a files_and_cmds")
	assert.Contains(t, out, `if [[ "$funcstack[1]" = "_desk" ]]; then`)
	assert.False(t, strings.Contains(out, "compadd -a matches"))
}
