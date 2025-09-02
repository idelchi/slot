# slot wrapper for zsh
slot() {
  emulate -L zsh
  set -o pipefail

  if [[ "$1" == "run" ]]; then
    shift

    local do_exec=0
    local -a passthru=()
    for arg in "$@"; do
      [[ "$arg" == "--exec" ]] && { do_exec=1; continue; }
      passthru+=("$arg")
    done

    local rendered rc
    rendered=$(command slot run "${passthru[@]}")
    rc=$?

    if (( rc != 0 )); then
      # If the app printed its error to STDOUT, surface it.
      [[ -n "$rendered" ]] && print -u2 -- "$rendered"
      return $rc
    fi

    [[ -z "$rendered" ]] && return 0

    if (( do_exec )); then
      eval "$rendered"         # no history push
    else
      print -z -- "$rendered" 2>/dev/null || print -r -- "$rendered"
    fi
    return $?
  fi

  command slot "$@"
}
