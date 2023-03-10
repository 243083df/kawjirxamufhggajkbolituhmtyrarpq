package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"io"
	"log"
	mrand "math/rand"
	"net"
)

const (
	saltSize   = 64
	answerSize = 64
	difficulty = 2
)

var quotes = []string{
	"To be sent greeting; not by com- mandment or constraint, but by revelation and the word of wisdom, showing forth the order and will of God in the temporal salvation of all saints in the last days—",
	"Given for a principle with promise, adapted to the capacity of the weak and the weakest of all saints, who are or can be called saints.",
	"Behold, verily, thus saith the Lord unto you: In consequence of evils and designs which do and will exist in the hearts of conspiring men in the last days, I have warned you, and forewarn you, by giving unto you this word of wisdom by revelation—",
	"That inasmuch as any man drinketh wine or strong drink among you, behold it is not good,  neither meet in the sight of your Father, only in assembling yourselves together to offer up your sacraments before him.6 And, behold, this should be wine,  yea, pure wine of the grape of the vine, of your own make.",
	"And, again, strong drinks are not for the belly, but for the washing of your bodies.",
	"And again, tobacco is not for the body, neither for the belly, and is not good for man, but is an herb for bruises and all sick cattle, to be used with judgment and skill.",
	"And again, hot drinks are not for the body or belly.",
	"And again, verily I say unto you,  all wholesome herbs God hath ordained for the constitution, nature,  and use of man—",
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Println("TCP server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err.Error())
			continue
		}

		challenge := append([]byte{difficulty, answerSize}, make([]byte, saltSize)...)
		salt := challenge[2:]

		_, err = io.ReadFull(crand.Reader, salt)
		if err != nil {
			log.Printf("Failed to generate salt: %s\n", err.Error())
			conn.Close()
			continue
		}

		_, err = conn.Write(challenge)
		if err != nil {
			log.Printf("Failed to send salt to client: %s\n", err.Error())
			conn.Close()
			continue
		}

		response := make([]byte, answerSize)
		_, err = conn.Read(response)
		if err != nil {
			log.Printf("Failed to read response from client: %s\n", err.Error())
			conn.Close()
			continue
		}

		if verifyChallenge(salt, response, difficulty) {
			quote := quotes[mrand.Intn(len(quotes))]
			_, err = conn.Write([]byte(quote))
			if err != nil {
				log.Printf("Failed to send quote to client: %s\n", err.Error())
			}
		}
		conn.Close()
	}
}

func verifyChallenge(salt, answer []byte, difficulty int) bool {
	hash := sha256.Sum256(append(salt, answer...))
	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}
