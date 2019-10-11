
log_error () {
  printf '\e[31mERROR: %s\n\e[39m' "$1" >&2
}
