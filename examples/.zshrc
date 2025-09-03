# slot integration
zmodload zsh/zle
# Ctrl-Y: run command in buffer as a slot
slot-run-buffer() {
  emulate -L zsh
  local buf="${BUFFER//$'\n'/ }"
  [[ -z "$buf" ]] && { zle -M "slot: buffer empty"; return 0; }

  BUFFER="slot run -y -- ${buf}"
  zle accept-line
}
zle -N slot-run-buffer
bindkey '^Z' slot-run-buffer


# Ctrl-X: show menu from `slot ls`
slot-pick-and-run() {
  emulate -L zsh
  set -o pipefail

  local choice name

  choice=$(
    slot ls \
    | head -n -1 \
    | fzf --prompt="slot> " --height=40% --header-lines=1 --layout=reverse-list
  ) || return

  name=${choice%%[[:space:]]*}
  BUFFER="slot run -y -- ${name}"
  zle accept-line
}
zle -N slot-pick-and-run
bindkey '^X' slot-pick-and-run
