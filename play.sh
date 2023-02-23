#!/bin/bash

# start a new game and get the game id
GAME_URL=$(curl -v POST -d '{"board":"---------"}' http://localhost:8080/api/v1/games 2>&1 | grep Location | awk '{print$3}' | tr -dc '[[:print:]]')

#Get the game id and make move
BOARD=$(curl --silent ${GAME_URL} | jq -r .board)

echo $BOARD

while true; do

	read -p "Enter index (0-8): " index
	read -p "Enter symbol (X): " symbol

	if [[ $index =~ ^[0-8]$ ]] && [[ $symbol =~ ^[XOxo]$ ]]; then
	   CLIENT_MOVE=${BOARD:0:$index}$symbol${BOARD:$((index+1))}
	fi

	echo Your move: $CLIENT_MOVE

	data="{\"board\":\"$CLIENT_MOVE\"}"

	json=$(curl --silent -X PUT -d $data ${GAME_URL})

	BOARD=$(echo $json | jq -r .board)

	echo Server move: $BOARD

	if [ $(echo $json | jq -r .status) != "RUNNING" ]; then
		echo $json | jq -r .status
		break
	fi
done
