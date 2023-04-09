package main

import (
	"fmt"
	"github.com/Jero075/GoMessenger-V2/encryption"
)

func main() {
	fmt.Println([]byte("Hello, World!!!!"))
	enc := encryption.GenerateCiphertext("test", "test", []byte("Hello, World!!!!"))
	fmt.Println(enc)
	dec := encryption.GeneratePlaintext("test", "test", enc)
	fmt.Println(dec)
}
