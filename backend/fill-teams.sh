#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if [ $# -ne 1 ]; then
  echo "Usage: $(basename "$0") <tournament-id>" >&2
  exit 1
fi

tournament_id="$1"
team_size=2

bridge-tab() {
  go run ./cli/main.go "$@"
}

# `tournament team list` prints a dashed rule before every team. A pattern that
# accepts hyphens treats those rules as team ids, joins a contestant, then fails
# with "no such team in Tournament".
listing=$(bridge-tab tournament team list -t "$tournament_id")

teams=$(printf '%s\n' "$listing" | awk '
  function uuid_count(line,    copy) {
    copy = line
    return gsub(/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/, "&", copy)
  }
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/ {
    if (team_id != "") print team_id, members
    team_id = substr($0, 1, 36)
    members = uuid_count($0) - 1
    next
  }
  team_id != "" { members += uuid_count($0) }
  END { if (team_id != "") print team_id, members }
')

if [ -z "$teams" ]; then
  echo "No teams found in tournament $tournament_id" >&2
  exit 1
fi

while read -r team_id member_count; do
  open_seats=$((team_size - member_count))
  if [ "$open_seats" -le 0 ]; then
    continue
  fi
  for ((seat = 0; seat < open_seats; seat++)); do
    contestant_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
    bridge-tab tournament join -i "$tournament_id" -c "$contestant_id"
    bridge-tab tournament team join -t "$tournament_id" -c "$contestant_id" -i "$team_id"
  done
done <<< "$teams"
