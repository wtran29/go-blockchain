package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := sign(); err != nil {
		log.Fatalln(err)
	}
}

func sign() error {
	// Need to load the private key file for the configured beneficiary so the
	// account can get credited with fees and tips.
	path := fmt.Sprintf("%s%s.ecdsa", "block/accounts/", "billy")
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println(address)

	v := struct {
		Name string
	}{
		Name: "Bill",
	}

	data, err := stamp(v)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}
	// Sign the hash with the private key to produce a signature
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	fmt.Println("Signature:", sig)
	fmt.Printf("SIG: 0x%s\n:", hex.EncodeToString(sig))

	// ===================================================================
	// NODE

	// Passed with sig
	v2 := struct {
		Name string
	}{
		Name: "Bill",
	}

	data2, err := stamp(v2)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}

	sigPublicKey, err := crypto.Ecrecover(data2, sig)
	if err != nil {
		return err
	}
	fmt.Println("PK_length:", len(sigPublicKey))

	// Extract the public key from the data and the signature.
	// publicKey, err := crypto.SigToPub(data2, sig)
	// if err != nil {
	// 	return err
	// }

	// Check the public key extracted from the data and signature.
	rs := sig[:crypto.RecoveryIDOffset]
	if !crypto.VerifySignature(sig, data2, rs) {
		return errors.New("invalid signature")
	}

	// Capture the public key associated with this data and signature.
	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	publicKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	// Extract the account address from the public key.
	address = crypto.PubkeyToAddress(publicKey).String()
	fmt.Println(address)

	return nil
}

func stamp(value any) ([]byte, error) {

	// Marshal the data.
	v, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// This stamp is used so signatures we produce when signing data
	// are always unique to the Go blockchain.
	stamp := []byte(fmt.Sprintf("\x19Signed Message:\n%d", len(v)))

	// Hash the stamp and txHash together in a final 32 byte array
	// that represents the data.
	data := crypto.Keccak256(stamp, v)

	return data, nil
}
