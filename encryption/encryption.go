package encryption

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"
)

func GenerateKeys() (publicKey rsa.PublicKey, privateKey crypto.PrivateKey) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating privateKey: " + err.Error())
	}

	pubKey := priKey.PublicKey
	return pubKey, priKey
}
func GenerateHash256(salt string, str string) (ret string) {
	hash := sha256.Sum256([]byte(salt + str + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func GenerateHash512(salt string, str string) (ret string) {
	hash := sha512.Sum512([]byte(salt + str + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func GenerateCiphertext(id string, password string, plaintext []byte) (ret []byte) {
	key := GenerateHash256(id, password)

	c, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating ciphertext: " + err.Error())
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating ciphertext: " + err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating ciphertext: " + err.Error())
	}
	return gcm.Seal(nonce, nonce, plaintext, nil)
}
func GeneratePlaintext(id string, password string, ciphertext []byte) (ret []byte) {
	key := GenerateHash256(id, password)

	c, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating plaintext: " + err.Error())
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating plaintext: " + err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating plaintext: " + err.Error())
	}
	return plaintext
}
func Encrypt(pubKey rsa.PublicKey, plaintext []byte) (ret []byte) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pubKey, plaintext, nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error encrypting message: " + err.Error())
	}

	return ciphertext
}
func Decrypt(priKey crypto.PrivateKey, ciphertext []byte) (ret []byte) {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priKey.(*rsa.PrivateKey), ciphertext, nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error decrypting message: " + err.Error())
	}

	return plaintext
}
