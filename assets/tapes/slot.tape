Output assets/gifs/slot.gif

Set Theme "iTerm2 Default"
Set FontSize 14
Set Width 1200
Set Height 800
Set TypingSpeed 75ms

Set Shell zsh
Env ZSH_AUTOSUGGEST_MANUAL_REBIND "1"
Env ZSH_AUTOSUGGEST_ENABLED "0"
Env SLOTS_FILE "~/assets/examples/slots.yaml"

# Setup
Hide
Type 'eval "$(slot init zsh --fzf)"'
Enter
Type 'setopt interactivecomments'
Enter
Type "clear"
Enter
Show

Type '# Add a slot'
Enter
Type 'slot save list "ls -la"'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# List all slots'
Enter
Type 'slot ls'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# Use fzf integration (ctrl-X)'
Enter
Ctrl+X
Sleep 2s
Type 'di'
Sleep 2s
Enter

Sleep 1s

# Finally
Sleep 4s

# Cleanup
Hide
Type 'slot rm list'
Enter
Type "clear"
Enter
