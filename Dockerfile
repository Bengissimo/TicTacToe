FROM golang:alpine AS builder

# create a working directory inside the image
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build  -o tictactoe cmd/main.go

FROM scratch

WORKDIR /build

COPY --from=builder /build/tictactoe /build/tictactoe

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8080

# command to be used to execute when the image is used to start a container
CMD ["./tictactoe"]