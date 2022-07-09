package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	c, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("Error during dial: ", err)
	}
	defer c.Close()

	connectionReader := bufio.NewReader(c)

	// Get nonce
	// Compute problem
	// Send hash
	// Get a word of wisdom
	// Exit
	for {
		_, err := c.Write([]byte("nonce\n"))
		if err != nil {
			log.Fatal("Error during writing to connection: ", err)
		}
		serverOutput, err := connectionReader.ReadString('\n')
		if err != nil {
			log.Fatal("Error during reading from connection: ", err)
		}
		fmt.Println(serverOutput)

		nonceAndDifficulty := strings.Split(serverOutput, ",")
		nonce, _ := strconv.Atoi(nonceAndDifficulty[0])
		// difficulty couldn't be more than 32
		difficulty, _ := strconv.Atoi(strings.TrimSpace(nonceAndDifficulty[1]))

		solutionNonce := 1
		fmt.Println(nonce, difficulty, solutionNonce)

		var hash [32]byte
		// PoW algorithm is here
		for {
			str := fmt.Sprintf("%d,%d", nonce, solutionNonce)
			hash = sha256.Sum256([]byte(str))

			zeroesToSolve := difficulty
			for i := 0; zeroesToSolve > 0; i++ {
				if hash[i] != 0 {
					break
				}
				zeroesToSolve--
			}

			if zeroesToSolve == 0 {
				fmt.Println(solutionNonce, string(hash[:]))
				break
			}

			solutionNonce++
		}

		// fmt.Printf("%#v, %d\n", hash, solutionNonce)
		_, err = c.Write([]byte(fmt.Sprintf("solve,%d\n", solutionNonce)))
		if err != nil {
			log.Fatal("Error during writing to connection: ", err)
		}

		_, err = c.Write([]byte(fmt.Sprintf("wisdom\n")))

		serverOutput, err = connectionReader.ReadString('\n')
		if err != nil {
			log.Fatal("Error during reading from connection: ", err)
		}

		return
	}
}
