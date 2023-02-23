# TicTacToe

A REST backend for a  Tic-tac-toe game. I used Go and the <a href="https://gin-gonic.com/docs/">Gin Web Framework</a> to create this RESTful web service API. 

The game board grid looks as follows
```
.-----------.
| 0 | 1 | 2 |
+---+---+---+
| 3 | 4 | 5 |
+---+---+---+
| 6 | 7 | 8 |
`-----------´
```

So, a board position
```
.-----------.
| X | O | - |
+---+---+---+
| - | X | - |
+---+---+---+
| - | O | X |
`-----------´
```

translates to
```
XO--X--OX
012345678
```
## Game flow:
----------

- Start a game with either empty (server starts) or with the first move made (client starts). The backend responds with the location URL of the started game:
```
$ curl -v POST -d '{"board":"--X------"}' http://localhost:8080/api/v1/games
...
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: http://127.0.0.1:8080/api/v1/games/3667fb47-fc9a-493a-8da6-a4190275bd20
< Date: Thu, 23 Feb 2023 13:50:45 GMT
< Content-Length: 84
< 
...
{"id":"3667fb47-fc9a-493a-8da6-a4190275bd20","board":"-OX------","status":"RUNNING"}
```

- Client GETs the board state from the URL:
```
$ curl http://127.0.0.1:8080/api/v1/games/3667fb47-fc9a-493a-8da6-a4190275bd20
{"id":"3667fb47-fc9a-493a-8da6-a4190275bd20","board":"-OX------","status":"RUNNING"}
```

- Client PUTs the board state with a new move to the URL. Backend validates the move, makes it's own move and updates the game state. The updated game state is returned in the PUT response:
```
$ curl -X PUT -d '{"board":"-OXX-----"}' http://127.0.0.1:8080/api/v1/games/3667fb47-fc9a-493a-8da6-a4190275bd20
{"id":"3667fb47-fc9a-493a-8da6-a4190275bd20","board":"-OXX--O--","status":"RUNNING"}
```

- And so on. The game is over once the computer or the player gets 3 noughts
  or crosses, horizontally, vertically or diagonally or there are no moves to
  be made:
```
$ curl -X PUT -d '{"board":"XOXXXOOOX"}' http://127.0.0.1:8080/api/v1/games/3667fb47-fc9a-493a-8da6-a4190275bd20
{"id":"3667fb47-fc9a-493a-8da6-a4190275bd20","board":"XOXXXOOOX","status":"X_WON"}% 
```

## Prerequisites
- Golang version 1.19 or higher
- `make`
- `jq` (to run play.sh)

## Usage
```
make all
```
Or you can run a docker image:
```
make docker
```
## Testing
### Unit tests
```
make test
```
### Play testing
```
./play.sh
```