package main

import (
	"crypto/sha256"
	"log"
	"math/rand"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("Failed to connect to server: %s\n", err.Error())
		return
	}

	challenge := make([]byte, 64+2)
	_, err = conn.Read(challenge)
	if err != nil {
		log.Printf("Failed to read challenge from server: %s\n", err.Error())
		return
	}

	difficulty := challenge[0]
	answerSize := challenge[1]
	salt := challenge[2:]
	answer, err := solveChallenge(salt, difficulty, answerSize)
	if err != nil {
		log.Printf("Failed to solve challenge: %s\n", err.Error())
		return
	}

	_, err = conn.Write(answer)
	if err != nil {
		log.Printf("Failed to send answer to server: %s\n", err.Error())
		return
	}

	quote := make([]byte, 1024)
	_, err = conn.Read(quote)
	if err != nil {
		log.Printf("Failed to read quote from server: %s\n", err.Error())
		return
	}

	log.Printf("Quote: %s\n", string(quote))
}

func solveChallenge(challenge []byte, difficulty byte, answerSize byte) ([]byte, error) {
	var answer []byte

	for {
		answer = make([]byte, answerSize)
		_, err := rand.Read(answer)
		if err != nil {
			return nil, err
		}

		if verifyChallenge(challenge, answer, difficulty) {
			return answer, nil
		}
	}
}

func verifyChallenge(salt, answer []byte, difficulty byte) bool {
	hash := sha256.Sum256(append(salt, answer...))
	for i := byte(0); i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}
