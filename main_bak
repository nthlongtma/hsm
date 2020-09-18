package main

import (
	"encoding/base64"
	"log"
)

const (
	configPath = "./configs"
	configName = "dev"
)

func main() {
	// read config
	// conf := LoadConfig(configPath, configName)
	// if conf == nil {
	// 	panic("failed to load config")
	// }
	// fmt.Printf("config: %+v\n", *conf)

	// XOR()

	// log.Println("simple encryption and decryption - XOR cipher")
	plain := []byte("This is a secret message.")
	key := []byte("secret")

	cipher := SimpleEncrypt(plain, key)
	log.Printf("cipher: %s", base64.StdEncoding.EncodeToString(cipher))

	decrypted := SimpleDecrypt(cipher, key)
	log.Printf("decrypted: %s", decrypted)
}

// XOR example
func XOR() {
	p, k := 19, 20

	log.Printf("k  : %08b = %[1]v", k)
	log.Printf("p  : %08b = %[1]v", p)
	// log.Printf("p^p: %08b = %[1]v", p^p)
	// log.Printf("p^0: %08b = %[1]v", p^0)

	c := p ^ k
	log.Printf("c  : %08b = %[1]v", c)
	log.Printf("d  : %08b = %[1]v", c^k) // c ^ k = p ^ k ^ K = p ^ 0 = p
}

func SimpleEncrypt(plain, key []byte) []byte {
	cipher := make([]byte, 0)

	if len(key) == 0 {
		return nil
	}

	for i := range plain {
		c := plain[i] ^ key[i%len(key)]
		cipher = append(cipher, c)
	}

	return cipher
}

func SimpleDecrypt(cipher, key []byte) []byte {
	decrypted := make([]byte, 0)

	if len(key) == 0 {
		return nil
	}

	for i := range cipher {
		d := cipher[i] ^ key[i%len(key)]
		decrypted = append(decrypted, d)
	}

	return decrypted
}
