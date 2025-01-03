alias bridge-tab='go run ./cli/main.go'

# Create tournament
message=$(bridge-tab tournament create -n "TEST - IX Turniej w Zagorzynie")
echo $message
tournament_id=$(echo $message | grep -oE '[0-9a-f-]{36}')

# Create teams
team_ids=()
for i in {1..12}; do
    message=$(bridge-tab tournament team create -t $tournament_id -n "$i")
    echo $message
    team_id=$(echo $message | grep -oE '[0-9a-f-]{36}' | head -1)
    team_ids+=($team_id)
done

# Create board protocols
bridge-tab tournament board-protocol create -i $tournament_id -n "1" -v "None" "12;1" "11;9" "5;8" "7;2" "6;10" "3;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "2" -v "NS" "12;1" "11;9" "5;8" "7;2" "6;10" "3;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "3" -v "EW" "12;1" "11;9" "5;8" "7;2" "6;10" "3;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "4" -v "Both" "12;2" "1;10" "6;9" "8;3" "7;11" "4;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "5" -v "NS" "12;2" "1;10" "6;9" "8;3" "7;11" "4;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "6" -v "EW" "12;2" "1;10" "6;9" "8;3" "7;11" "4;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "7" -v "Both" "12;3" "2;11" "7;10" "9;4" "8;1" "5;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "8" -v "None" "12;3" "2;11" "7;10" "9;4" "8;1" "5;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "9" -v "EW" "12;3" "2;11" "7;10" "9;4" "8;1" "5;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "10" -v "Both" "6;7" "12;4" "3;1" "8;11" "10;5" "9;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "11" -v "None" "6;7" "12;4" "3;1" "8;11" "10;5" "9;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "12" -v "NS" "6;7" "12;4" "3;1" "8;11" "10;5" "9;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "13" -v "Both" "7;8" "12;5" "4;2" "9;1" "11;6" "10;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "14" -v "None" "7;8" "12;5" "4;2" "9;1" "11;6" "10;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "15" -v "NS" "7;8" "12;5" "4;2" "9;1" "11;6" "10;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "16" -v "EW" "11;4" "8;9" "12;6" "5;3" "10;2" "1;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "17" -v "None" "11;4" "8;9" "12;6" "5;3" "10;2" "1;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "18" -v "NS" "11;4" "8;9" "12;6" "5;3" "10;2" "1;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "19" -v "EW" "1;5" "9;10" "12;7" "6;4" "11;3" "2;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "20" -v "Both" "1;5" "9;10" "12;7" "6;4" "11;3" "2;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "21" -v "NS" "1;5" "9;10" "12;7" "6;4" "11;3" "2;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "22" -v "EW" "3;9" "2;6" "10;11" "12;8" "7;5" "1;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "23" -v "Both" "3;9" "2;6" "10;11" "12;8" "7;5" "1;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "24" -v "None" "3;9" "2;6" "10;11" "12;8" "7;5" "1;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "25" -v "EW" "2;5" "4;10" "3;7" "11;1" "12;9" "8;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "26" -v "Both" "2;5" "4;10" "3;7" "11;1" "12;9" "8;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "27" -v "None" "2;5" "4;10" "3;7" "11;1" "12;9" "8;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "28" -v "NS" "3;6" "5;11" "4;8" "1;2" "12;10" "9;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "29" -v "Both" "3;6" "5;11" "4;8" "1;2" "12;10" "9;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "30" -v "None" "3;6" "5;11" "4;8" "1;2" "12;10" "9;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "31" -v "NS" "10;8" "4;7" "6;1" "5;9" "2;3" "12;11"
bridge-tab tournament board-protocol create -i $tournament_id -n "32" -v "EW" "10;8" "4;7" "6;1" "5;9" "2;3" "12;11"
bridge-tab tournament board-protocol create -i $tournament_id -n "33" -v "None" "10;8" "4;7" "6;1" "5;9" "2;3" "12;11"

# Create contestants (testing purposes)
for team_id in "${team_ids[@]}"; do
    uuid=$(uuidgen)
    bridge-tab tournament join -i $tournament_id -c $uuid
    echo $team_id
    bridge-tab tournament team join -t $tournament_id -c $uuid -i $team_id
done
