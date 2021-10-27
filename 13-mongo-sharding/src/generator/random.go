package generator

import "math/rand"

var Alphabet = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")

func GetRandomKey() string {
	idBytes := make([]byte, 5)
	for i := 0; i < len(idBytes); i++ {
		idBytes[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return string(idBytes)
}
