package encryption

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

func GenerateKeys(username string, password string) (publicKey ecdsa.PublicKey, privateKey ecdsa.PrivateKey) {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating privateKey: " + err.Error())
	}

	pubKey := priKey.PublicKey
	priKeyEnc := new(big.Int).SetBytes(GenerateCiphertext(username, password, priKey.D.Bytes()))
	return pubKey, ecdsa.PrivateKey{priKey.PublicKey, priKeyEnc}
}
func GenerateHash256(salt string, str string) (ret string) {
	hash := sha256.Sum256([]byte(salt + str + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func GenerateHash512(salt string, str string) (ret string) {
	hash := sha512.Sum512([]byte(salt + str + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func GenerateCiphertext(username string, password string, plaintext []byte) (ret []byte) {
	key16 := GenerateHash256(username, password)

	c, err := aes.NewCipher([]byte(key16[:16]))
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating ciphertext: " + err.Error())
	}
	ct := make([]byte, len(plaintext))
	c.Encrypt(ct, plaintext)

	return ct
}
func GeneratePlaintext(username string, password string, ciphertext []byte) (ret []byte) {
	key16 := GenerateHash256(username, password)

	c, err := aes.NewCipher([]byte(key16[:16]))
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating plaintext: " + err.Error())
	}
	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return pt
}
