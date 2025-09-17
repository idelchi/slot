# slot key-bindings for Ctrl-X (using fzf) and Ctrl-Z

__slot_eval_prompt() {
  if ((BASH_VERSINFO[0] > 4 || (BASH_VERSINFO[0] == 4 && BASH_VERSINFO[1] >= 4))); then
    local p=${PS1@P}
    printf '%s' "${p//$'\001'/}" | tr -d $'\002'
  else
    printf '%s' "$PS1"
  fi
}
__slot_accept_line() {
  local __cmd=$1
  local __p __ret __stty
  __p="$(__slot_eval_prompt)"
  printf '\r'; tput el 2>/dev/null || true
  printf '%s%s\n' "$__p" "$__cmd"
  builtin history -s "$__cmd"
  __stty=$(stty -g 2>/dev/null || true)
  eval -- "$__cmd"$'\n__ret=$?'
  [[ -n $__stty ]] && stty "$__stty" 2>/dev/null || true
  if [[ -n ${PROMPT_COMMAND-} ]]; then
    if declare -p PROMPT_COMMAND 2>/dev/null | grep -q 'declare \-a'; then
      local __pc
      for __pc in "${PROMPT_COMMAND[@]}"; do eval -- "$__pc"; done
    else
      eval -- "$PROMPT_COMMAND"
    fi
  fi
  return $__ret
}

# Ctrl-Y: run command in buffer as a slot
slot_run_buffer() {
  local buf=${READLINE_LINE//$'\n'/ }
  if [[ -z "$buf" ]]; then
    echo "no slot selected"
    return 0
  fi
  __slot_accept_line "slot run -y -- ${buf}"
  READLINE_LINE=
  READLINE_POINT=0
}

# Ctrl-X: show menu from `slot ls --tsv`
slot_pick_and_run() {
  set -o pipefail
  local out key choice
  local name tags cmd

  out=$(
    slot ls --tsv | fzf \
      --prompt="slot> " \
      --height=40% \
      --layout=reverse-list \
      --header $'ENTER: run  TAB: insert slot  SHIFT-TAB: insert CMD' \
      --header-lines=1 \
      --delimiter=$'\t' \
      --with-nth=1,2,3 \
      --tabstop=16 \
      --expect=enter,tab,btab
  ) || return

  key=${out%%$'\n'*}
  choice=${out#*$'\n'}
  [[ -z $choice || $choice = "$key" ]] && return

  IFS=$'\t' read -r name tags cmd <<<"$choice"
  cmd=${cmd//\\n/$'\n'}

  case $key in
    enter) __slot_accept_line "slot run -y -- ${name}"; READLINE_LINE=; READLINE_POINT=0; return ;;
    tab)   READLINE_LINE="slot run -y -- ${name}" ;;
    btab)  READLINE_LINE="${cmd}" ;;
  esac

  READLINE_POINT=${#READLINE_LINE}
}

stty susp undef 2>/dev/null || true
bind -r '\C-z' 2>/dev/null || true
bind -x '"\C-z": slot_run_buffer'
bind -x '"\C-x": slot_pick_and_run'
