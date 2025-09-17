# slot wrapper for zsh
slot() {
  emulate -L zsh
  set -o pipefail

  if [[ "$1" == "run" ]]; then
    shift

    local do_exec=0
    local -a passthru=()
    for arg in "$@"; do
      if [[ "${arg}" == "--yes" || "${arg}" == "-y" ]]; then
        do_exec=1
        continue
      fi
      passthru+=("${arg}")
    done

    local rendered rc
    rendered=$(command slot render "${passthru[@]}")
    rc=$?

    if (( rc != 0 )); then
      [[ -n "${rendered}" ]] && print -u2 -- "${rendered}"
      return ${rc}
    fi

    [[ -z "${rendered}" ]] && return 0

    if (( do_exec )); then
      eval "${rendered}"
    else
      print -z -- "${rendered}" 2>/dev/null || print -r -- "${rendered}"
    fi
    return $?
  fi

  command slot "$@"
}
