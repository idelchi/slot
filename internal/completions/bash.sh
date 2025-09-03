# slot wrapper for bash
slot() {
  set -o pipefail

  if [[ "$1" == "run" ]]; then
    shift

    local do_exec=0
    local -a passthru=()
    local arg
    for arg in "$@"; do
      if [[ "${arg}" == "--yes" || "${arg}" == "-y" ]]; then
        do_exec=1
      else
        passthru+=("${arg}")
      fi
    done

    # capture stdout from the real 'slot' command
    local rendered rc
    rendered="$(command slot render "${passthru[@]}")"
    rc=$?

    if (( rc != 0 )); then
      # surface any error text the tool printed
      if [[ -n "${rendered}" ]]; then
        printf '%s\n' "${rendered}" >&2
      fi
      return "$rc"
    fi

    # nothing to do
    [[ -z "${rendered}" ]] && return 0

    if (( do_exec )); then
      eval "${rendered}"   # no history push
    else
      # Bash has no 'print -z'. Best approximation:
      #   1) push to history so it's available with Up-arrow
      #   2) also echo it so you can see/copy it immediately
      history -s "${rendered}"
      printf '%s\n' "${rendered}"
    fi
    return $?
  fi

  command slot "$@"
}
