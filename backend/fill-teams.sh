#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if [ $# -ne 1 ]; then
  echo "Usage: $(basename "$0") <tournament-id>" >&2
  exit 1
fi

tournament_id="$1"

bridge-tab() {
  go run ./cli/main.go "$@"
}

team_ids=$(bridge-tab tournament team list -t "$tournament_id" | awk '/^[0-9a-f-]{36}/ { print substr($0, 1, 36) }')

if [ -z "$team_ids" ]; then
  echo "No teams found in tournament $tournament_id" >&2
  exit 1
fi

while IFS= read -r team_id; do
  contestant_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
  bridge-tab tournament join -i "$tournament_id" -c "$contestant_id"
  bridge-tab tournament team join -t "$tournament_id" -c "$contestant_id" -i "$team_id"
done <<< "$team_ids"
