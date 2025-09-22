# slot key-bindings for Ctrl-X (using fzf) and Ctrl-Z
zmodload zsh/zle

# Ctrl-Y: run command in buffer as a slot
slot-run-buffer() {
  emulate -L zsh
  local buf="${BUFFER//$'\n'/ }"

  if [[ -z $buf ]]; then
    zle -I
    print -r -- "no slot selected"
    return 0
  fi

  BUFFER="slot run -y -- ${buf}"
  zle accept-line
}
zle -N slot-run-buffer
bindkey '^Z' slot-run-buffer

# Ctrl-X: show menu from `slot ls --tsv`
slot-pick-and-run() {
  emulate -L zsh
  set -o pipefail
  local out key choice
  local -a fields
  local name cmd

  out=$(
    slot ls --tsv | fzf \
      --prompt="slot> " \
      --height=40% \
      --layout=reverse-list \
      --header $'ENTER: run  TAB: insert slot  SHIFT-TAB: insert CMD  CTRL-SPACE: insert rendered cmd  CTRL-R: toggle preview' \
      --header-lines=1 \
      --delimiter=$'\t' \
      --nth=1,2,3,4 \
      --with-nth=1,3,4 \
      --tabstop=16 \
      --preview 'printf "%s\n" {2} | sed -e "s/\\^J/\n/g" -e "s/\\\\n/\n/g"' \
      --preview-window=25% \
      --bind 'ctrl-r:toggle-preview' \
      --style=full \
      --expect=enter,tab,btab,ctrl-space
  ) || { zle reset-prompt; return }

  key=${out%%$'\n'*}
  choice=${out#*$'\n'}
  [[ -z $choice || $choice = "$key" ]] && { zle reset-prompt; return }

  fields=("${(@ps:\t:)choice}")
  name=${fields[1]}
  cmd=${fields[2]//\\n/$'\n'}
  cmd=${cmd//^J/$'\n'}

  case $key in
    enter) BUFFER="slot run -y ${name}"; zle accept-line; return ;;
    tab)   BUFFER="slot run -y ${name}" ;;
    btab)  BUFFER="${cmd}" ;;
    ctrl-space) BUFFER="$(slot render ${name})" ;;
  esac

  CURSOR=${#BUFFER}
  zle redisplay
}
zle -N slot-pick-and-run
bindkey '^X' slot-pick-and-run
