package sign

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// DocumentHash represents the hashed value of a document
type DocumentHash string

// Signature represents a document's signature
type Signature string

// The types DocumentHash and Signature are implemented as separate types to prevent accidental mixing up of these values
// which was possible when the functions creating them simply returned strings instead of these custom stirng based types

// CreateKeyPair creates and returns a private key, public key
func CreateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	privKey.Precompute()

	if err := privKey.Validate(); err != nil {
		return nil, nil, err
	}

	pubKey := privKey.PublicKey
	return privKey, &pubKey, nil
}

// HashDocument creates and returns the DocumentHash of the provided document (d)
func HashDocument(d interface{}) (DocumentHash, error) {
	var buf = bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, d)
	return DocumentHash(fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))), nil
}

// Sign creates and returns a Signature based upon a document hash (h) and a private key (k)
func Sign(h DocumentHash, k *rsa.PrivateKey) (Signature, error) {
	hash, err := hex.DecodeString(string(h))
	if err != nil {
		return "", err
	}

	sig, err := rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return Signature(fmt.Sprintf("%x", sig)), nil
}

// ValidSignature checks if the provided public key (k) is valid for the signature (s) of a document's hash (h)
func ValidSignature(h DocumentHash, s Signature, k *rsa.PublicKey) error {
	hash, err := hex.DecodeString(string(h))
	if err != nil {
		return err
	}

	sig, err := hex.DecodeString(string(s))
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(k, crypto.SHA256, hash[:], sig)
	if err != nil {
		return err
	}
	return nil
}

// Todo: Examples
