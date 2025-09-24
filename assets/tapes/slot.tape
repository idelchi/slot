Output assets/gifs/slot.gif

Set Theme "iTerm2 Default"
Set FontSize 14
Set Width 1200
Set Height 800
Set TypingSpeed 75ms

Set Shell bash
Env SLOTS_FILE "~/assets/examples/slots.yaml"

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
Sleep 500ms
Type 'eval "$(slot init bash --fzf)"'
Sleep 500ms
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
