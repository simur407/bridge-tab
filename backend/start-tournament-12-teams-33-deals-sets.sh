alias bridge-tab='go run ./cli/main.go'

# Create tournament
message=$(bridge-tab tournament create -n "Demonstracyjny")
echo $message
tournament_id=$(echo $message | grep -oE '[0-9a-f-]{36}')

# Create teams
team_ids=()
for i in {1..12}; do
    message=$(bridge-tab tournament team create -t $tournament_id --number "$i")
    echo $message
    team_id=$(echo $message | grep -oE '[0-9a-f-]{36}' | head -1)
    team_ids+=($team_id)
done

# Create tables
for i in {1..6}; do
    bridge-tab tournament table create -t $tournament_id --number "$i"
done

# Create board protocols.
# Pairings are "{table}:{NS};{EW}" from the 6-table Howell.
# Relay boards are not in play. Each set is played at the six tables below.
# A 1:12-1  2:9-4  3:3-10  4:7-6  5:11-2  6:8-5
# B 1:12-2  2:10-5 3:4-11  4:8-7  5:1-3   6:9-6
# C 1:12-3  2:11-6 3:5-1   4:9-8  5:2-4   6:10-7
# D 1:12-4  2:1-7  3:6-2   4:10-9 5:3-5   6:11-8
# E 1:12-5  2:2-8  3:7-3   4:11-10 5:4-6  6:1-9
# F 1:12-6  2:3-9  3:8-4   4:1-11 5:5-7   6:2-10
# G 1:12-7  2:4-10 3:9-5   4:2-1  5:6-8   6:3-11
# H 1:12-8  2:5-11 3:10-6  4:3-2  5:7-9   6:4-1
# I 1:12-9  2:6-1  3:11-7  4:4-3  5:8-10  6:5-2
# J 1:12-10 2:7-2  3:1-8   4:5-4  5:9-11  6:6-3
# K 1:12-11 2:8-3  3:2-9   4:6-5  5:10-1  6:7-4
bridge-tab tournament board-protocol create -i $tournament_id -n "1" -v "None" "1:12;1" "2:9;4" "3:3;10" "4:7;6" "5:11;2" "6:8;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "2" -v "NS" "1:12;1" "2:9;4" "3:3;10" "4:7;6" "5:11;2" "6:8;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "3" -v "EW" "1:12;1" "2:9;4" "3:3;10" "4:7;6" "5:11;2" "6:8;5"
bridge-tab tournament board-protocol create -i $tournament_id -n "4" -v "Both" "1:12;2" "2:10;5" "3:4;11" "4:8;7" "5:1;3" "6:9;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "5" -v "NS" "1:12;2" "2:10;5" "3:4;11" "4:8;7" "5:1;3" "6:9;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "6" -v "EW" "1:12;2" "2:10;5" "3:4;11" "4:8;7" "5:1;3" "6:9;6"
bridge-tab tournament board-protocol create -i $tournament_id -n "7" -v "Both" "1:12;3" "2:11;6" "3:5;1" "4:9;8" "5:2;4" "6:10;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "8" -v "None" "1:12;3" "2:11;6" "3:5;1" "4:9;8" "5:2;4" "6:10;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "9" -v "EW" "1:12;3" "2:11;6" "3:5;1" "4:9;8" "5:2;4" "6:10;7"
bridge-tab tournament board-protocol create -i $tournament_id -n "10" -v "Both" "1:12;4" "2:1;7" "3:6;2" "4:10;9" "5:3;5" "6:11;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "11" -v "None" "1:12;4" "2:1;7" "3:6;2" "4:10;9" "5:3;5" "6:11;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "12" -v "NS" "1:12;4" "2:1;7" "3:6;2" "4:10;9" "5:3;5" "6:11;8"
bridge-tab tournament board-protocol create -i $tournament_id -n "13" -v "Both" "1:12;5" "2:2;8" "3:7;3" "4:11;10" "5:4;6" "6:1;9"
bridge-tab tournament board-protocol create -i $tournament_id -n "14" -v "None" "1:12;5" "2:2;8" "3:7;3" "4:11;10" "5:4;6" "6:1;9"
bridge-tab tournament board-protocol create -i $tournament_id -n "15" -v "NS" "1:12;5" "2:2;8" "3:7;3" "4:11;10" "5:4;6" "6:1;9"
bridge-tab tournament board-protocol create -i $tournament_id -n "16" -v "EW" "1:12;6" "2:3;9" "3:8;4" "4:1;11" "5:5;7" "6:2;10"
bridge-tab tournament board-protocol create -i $tournament_id -n "17" -v "None" "1:12;6" "2:3;9" "3:8;4" "4:1;11" "5:5;7" "6:2;10"
bridge-tab tournament board-protocol create -i $tournament_id -n "18" -v "NS" "1:12;6" "2:3;9" "3:8;4" "4:1;11" "5:5;7" "6:2;10"
bridge-tab tournament board-protocol create -i $tournament_id -n "19" -v "EW" "1:12;7" "2:4;10" "3:9;5" "4:2;1" "5:6;8" "6:3;11"
bridge-tab tournament board-protocol create -i $tournament_id -n "20" -v "Both" "1:12;7" "2:4;10" "3:9;5" "4:2;1" "5:6;8" "6:3;11"
bridge-tab tournament board-protocol create -i $tournament_id -n "21" -v "NS" "1:12;7" "2:4;10" "3:9;5" "4:2;1" "5:6;8" "6:3;11"
bridge-tab tournament board-protocol create -i $tournament_id -n "22" -v "EW" "1:12;8" "2:5;11" "3:10;6" "4:3;2" "5:7;9" "6:4;1"
bridge-tab tournament board-protocol create -i $tournament_id -n "23" -v "Both" "1:12;8" "2:5;11" "3:10;6" "4:3;2" "5:7;9" "6:4;1"
bridge-tab tournament board-protocol create -i $tournament_id -n "24" -v "None" "1:12;8" "2:5;11" "3:10;6" "4:3;2" "5:7;9" "6:4;1"
bridge-tab tournament board-protocol create -i $tournament_id -n "25" -v "EW" "1:12;9" "2:6;1" "3:11;7" "4:4;3" "5:8;10" "6:5;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "26" -v "Both" "1:12;9" "2:6;1" "3:11;7" "4:4;3" "5:8;10" "6:5;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "27" -v "None" "1:12;9" "2:6;1" "3:11;7" "4:4;3" "5:8;10" "6:5;2"
bridge-tab tournament board-protocol create -i $tournament_id -n "28" -v "NS" "1:12;10" "2:7;2" "3:1;8" "4:5;4" "5:9;11" "6:6;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "29" -v "Both" "1:12;10" "2:7;2" "3:1;8" "4:5;4" "5:9;11" "6:6;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "30" -v "None" "1:12;10" "2:7;2" "3:1;8" "4:5;4" "5:9;11" "6:6;3"
bridge-tab tournament board-protocol create -i $tournament_id -n "31" -v "NS" "1:12;11" "2:8;3" "3:2;9" "4:6;5" "5:10;1" "6:7;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "32" -v "EW" "1:12;11" "2:8;3" "3:2;9" "4:6;5" "5:10;1" "6:7;4"
bridge-tab tournament board-protocol create -i $tournament_id -n "33" -v "None" "1:12;11" "2:8;3" "3:2;9" "4:6;5" "5:10;1" "6:7;4"

# Create sets (boards that share the same team pairings)
bridge-tab tournament set create -i $tournament_id -l "A" -b "1,2,3"
bridge-tab tournament set create -i $tournament_id -l "B" -b "4,5,6"
bridge-tab tournament set create -i $tournament_id -l "C" -b "7,8,9"
bridge-tab tournament set create -i $tournament_id -l "D" -b "10,11,12"
bridge-tab tournament set create -i $tournament_id -l "E" -b "13,14,15"
bridge-tab tournament set create -i $tournament_id -l "F" -b "16,17,18"
bridge-tab tournament set create -i $tournament_id -l "G" -b "19,20,21"
bridge-tab tournament set create -i $tournament_id -l "H" -b "22,23,24"
bridge-tab tournament set create -i $tournament_id -l "I" -b "25,26,27"
bridge-tab tournament set create -i $tournament_id -l "J" -b "28,29,30"
bridge-tab tournament set create -i $tournament_id -l "K" -b "31,32,33"

# Create contestants (testing purposes)
for team_id in "${team_ids[@]}"; do
    uuid=$(uuidgen | tr '[:upper:]' '[:lower:]')
    bridge-tab tournament join -i $tournament_id -c $uuid
    echo $team_id
    bridge-tab tournament team join -t $tournament_id -c $uuid -i $team_id
done
