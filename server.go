package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
)

// Player description
type Player struct {
	ID          string
	IsConnected bool
	Score       int
	IsMoveMade  bool
	Choice      string
}

var player1, player2 *Player
var socketsServer *socketio.Server

func main() {
	player1 = new(Player)
	player2 = new(Player)
	localAddress := getLocalAddress()

	startSocketsServer()
	defer socketsServer.Close()

	http.Handle("/socket.io/", socketsServer)
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", fileServer)

	fmt.Println("Server Address is: " + localAddress.String() + ":5000")
	http.ListenAndServe(":5000", nil)
}

func startSocketsServer() {
	socketsServer = socketio.NewServer(nil)

	socketsServer.OnConnect("/", func(connection socketio.Conn) error {
		playerConnect(connection)
		return nil
	})

	socketsServer.OnDisconnect("/", func(connection socketio.Conn, s string) {
		playerDisconnect(connection)
	})

	socketsServer.OnEvent("/", "play", func(connection socketio.Conn, msg string) {
		play(connection, msg)
	})

	go socketsServer.Serve()
}

func play(connection socketio.Conn, msg string) {
	currentPlayer := getPlayerByID(connection.ID())

	currentPlayer.IsMoveMade = true
	currentPlayer.Choice = msg

	if player1.IsMoveMade && player2.IsMoveMade {
		makeMove()
	}
}

func makeMove() {
	const PaperMove = "paper"
	const ScissorsMove = "scissors"
	const RockMove = "rock"
	const P1WinsMessage = "P1 Wins"
	const P2WinsMessage = "P2 Wins"

	response := player1.Choice + "#" + player2.Choice
	message := ""

	if player1.Choice == player2.Choice {
		message = "Draw"
	} else if player1.Choice == PaperMove {
		if player2.Choice == ScissorsMove {
			player2.Score++
			message = P2WinsMessage
		} else {
			player1.Score++
			message = P1WinsMessage
		}
	} else if player1.Choice == RockMove {
		if player2.Choice == PaperMove {
			player2.Score++
			message = P2WinsMessage
		} else {
			player1.Score++
			message = P1WinsMessage
		}
	} else if player1.Choice == ScissorsMove {
		if player2.Choice == RockMove {
			player2.Score++
			message = P2WinsMessage
		} else {
			player1.Score++
			message = P1WinsMessage
		}
	}

	response = response + "#" + strconv.Itoa(player1.Score) + "#" +
		strconv.Itoa(player2.Score) + "#" + message

	player1.IsMoveMade = false
	player2.IsMoveMade = false

	socketsServer.BroadcastToRoom("", "gameRoom", "makeMove", response)
}

func playerConnect(connection socketio.Conn) {
	fmt.Println("Connected client: " + connection.ID())
	selectedPlayer := ""

	if player1.IsConnected == false {
		player1.IsConnected = true
		player1.ID = connection.ID()
		selectedPlayer = "P1"
		connection.Join("gameRoom")
	} else if player2.IsConnected == false {
		player2.IsConnected = true
		player2.ID = connection.ID()
		selectedPlayer = "P2"
		connection.Join("gameRoom")
	}

	connection.Emit("ready", selectedPlayer)
	if player1.IsConnected && player2.IsConnected && selectedPlayer != "" {
		socketsServer.BroadcastToRoom("", "gameRoom", "gameStart")
	}
}

func playerDisconnect(connection socketio.Conn) {
	connectionID := connection.ID()

	player := getPlayerByID(connectionID)
	if player != nil {
		player.IsConnected = false
		player.ID = ""
		resetGame(connection)
	}
}

func getPlayerByID(connectionID string) *Player {
	if connectionID == player1.ID {
		return player1
	} else if connectionID == player2.ID {
		return player2
	}

	return nil
}

func resetGame(connection socketio.Conn) {
	fmt.Println("Game RESET")
	player1 = new(Player)
	player2 = new(Player)
	socketsServer.BroadcastToRoom("", "gameRoom", "gameReset", "")
}

func getLocalAddress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	conn.Close()

	return localAddr.IP
}
