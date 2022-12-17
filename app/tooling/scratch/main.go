package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println(address)

	return nil
}
