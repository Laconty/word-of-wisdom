package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var difficulty = 3
var wordsOfWisdom []string = []string{
	"1",
	"31",
	"21",
	"12",
	"312",
	"212",
	"13",
	"313",
	"213",
	"14",
	"314",
	"214",
	"15",
	"315",
	"215",
	"16",
	"316",
	"216",
	"17",
	"317",
	"217",
	"18",
	"318",
	"218",
	"19",
	"319",
	"219",
}

func main() {
	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Printf("listening to port %s\n", listener.Addr().String())

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error acceoting new connection ", err)
		}

		go tcpHandle(c)
	}
}

var (
	CommandGetNonce     = "nonce"
	CommandSolveProblem = "solve"
	CommandGetWisdom    = "wisdom"
)

type ConnectionState struct {
	Nonce      int // don't confuse this Nonce and user's solved nonce
	Difficulty int
	Solved     bool // Solved should be timestamp, and could be changed dynamically
}

func tcpHandle(c net.Conn) {
	defer c.Close()
	connReader := bufio.NewReader(c)
	connectionState := ConnectionState{
		Difficulty: difficulty, // TODO: Should incrment difficulty dynamicly/by-user-fails
	}

	for {
		msg, err := connReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error during reading from connection ", err)
			}
			return
		}
		msg = strings.TrimSpace(msg)
		commandAndMessage := strings.Split(msg, ",")
		fmt.Printf("%#v\n", commandAndMessage)
		if len(commandAndMessage) < 1 {
			fmt.Fprintln(c, "error: format should be $COMMAND[,$MSG]")
			continue
		}

		switch strings.TrimSpace(commandAndMessage[0]) {
		case CommandGetNonce:
			rand.Seed(time.Now().UnixNano())
			connectionState.Nonce = rand.Int()
			// fmt.Println(connectionState.Nonce, connectionState.Difficulty)

			fmt.Fprintf(c, "%d,%d\n", connectionState.Nonce, connectionState.Difficulty)
		case CommandSolveProblem:
			if len(commandAndMessage) < 2 {
				fmt.Print(c, "error: format should be $COMMAND,$MSG\n")
				break
			}

			hash := sha256.Sum256([]byte(fmt.Sprintf("%d,%s", connectionState.Nonce, commandAndMessage[1])))
			zeroesToSolve := connectionState.Difficulty
			for i := 0; zeroesToSolve > 0; i++ {
				if hash[i] != 0 {
					break
				}
				zeroesToSolve--
			}
			if zeroesToSolve > 0 {
				fmt.Fprint(c, "error: mistake in solution!\n")
				break
			}

			fmt.Fprint(c, "solved successfully")
			connectionState.Solved = true
			// fmt.Printf("%#v", connectionState) // DEBUG
		case CommandGetWisdom:
			if !connectionState.Solved {
				fmt.Fprint(c, "error: solve problem first\n")
				break
			}

			fmt.Fprint(c, getWisdomWord())
		default:
			fmt.Fprint(c, "error: invalid command\n")
		}
	}
}

func getWisdomWord() string {
	rand.Seed(time.Now().UnixNano())
	return wordsOfWisdom[rand.Intn(len(wordsOfWisdom))] + "\n"
}
