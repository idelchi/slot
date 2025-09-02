# slot wrapper for bash
slot() {
  if [[ "$1" == "run" ]]; then
    shift
    local rendered
    rendered=$(command slot run "$@") || return $?
    if [[ " $@ " == *" --exec "* ]]; then
      # execute directly
      eval "$rendered"
    else
      # replace current line in readline buffer
      READLINE_LINE="$rendered"
      READLINE_POINT=${#READLINE_LINE}
    fi
    return
  fi
  command slot "$@"
}
bind -x '"\C-o": true'   # ensure READLINE_LINE replacement works
