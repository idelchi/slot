Output gifs/slot.gif

Set Theme "iTerm2 Default"
Set FontSize 14
Set Width 800
Set Height 700
Set TypingSpeed 75ms

Set Shell zsh
Env ZSH_AUTOSUGGEST_MANUAL_REBIND "1"
Env ZSH_AUTOSUGGEST_ENABLED "0"

# Setup
Hide
Type 'mkdir -p ~/.config/slot'
Enter
Type 'cp ~/slots.yaml ~/.config/slot/slots.yaml'
Enter
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
# Down
Sleep 2s
Enter

Sleep 1s

# Finally
Sleep 4s
