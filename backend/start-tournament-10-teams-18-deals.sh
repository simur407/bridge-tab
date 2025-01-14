alias bridge-tab='go run ./cli/main.go'

# Create tournament
message=$(bridge-tab tournament create -n "Demo")
echo $message
tournament_id=$(echo $message | grep -oE '[0-9a-f-]{36}')

# Create teams
team_ids=()
for i in {1..10}; do
    message=$(bridge-tab tournament team create -t $tournament_id -n "$i")
    echo $message
    team_id=$(echo $message | grep -oE '[0-9a-f-]{36}' | head -1)
    team_ids+=($team_id)
done

# Create board protocols
bridge-tab tournament board-protocol create -i $tournament_id -n "1" -v "None" "10;1" "6;3" "5;9" "4;2" "8;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "2" -v "NS" "10;1" "6;3" "5;9" "4;2" "8;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "3" -v "EW" "9;8" "10;2" "7;4" "6;1" "5;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "4" -v "Both" "9;8" "10;2" "7;4" "6;1" "5;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "5" -v "NS" "6;4" "1;9" "10;3" "8;5" "7;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "6" -v "EW" "6;4" "1;9" "10;3" "8;5" "7;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "7" -v "Both" "7;5" "2;1" "10;4" "9;6" "8;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "8" -v "None" "7;5" "2;1" "10;4" "9;6" "8;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "9" -v "EW" "8;6" "3;2" "10;5" "1;7" "9;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "10" -v "Both" "8;6" "3;2" "10;5" "1;7" "9;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "11" -v "None" "9;7" "4;3" "10;6" "2;8" "1;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "12" -v "NS" "9;7" "4;3" "10;6" "2;8" "1;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "13" -v "Both" "1;8" "5;4" "10;7" "3;9" "2;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "14" -v "None" "1;8" "5;4" "10;7" "3;9" "2;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "15" -v "NS" "3;7" "2;9" "6;5" "10;8" "4;1"
bridge-tab tournament board-protocol create -i $tournament_id -n "16" -v "EW" "3;7" "2;9" "6;5" "10;8" "4;1"
bridge-tab tournament board-protocol create -i $tournament_id -n "17" -v "None" "5;2" "4;8" "3;1" "7;6" "10;9"
bridge-tab tournament board-protocol create -i $tournament_id -n "18" -v "NS" "5;2" "4;8" "3;1" "7;6" "10;9"

# Create contestants (testing purposes)
for team_id in "${team_ids[@]}"; do
    uuid=$(uuidgen)
    bridge-tab tournament join -i $tournament_id -c $uuid
    echo $team_id
    bridge-tab tournament team join -t $tournament_id -c $uuid -i $team_id
done
