NAME=tictactoe

all: run

${NAME}:
	go build -o ${NAME} main.go

run: ${NAME}
	./${NAME}

clean:
	rm ${NAME}

test:
	go test ./... -v

docker_build:
	docker build -t bengissimo/tictactoe .

docker_run: docker_build
	docker run --rm -p 8080:8080 bengissimo/tictactoe

docker: docker_build docker_run
