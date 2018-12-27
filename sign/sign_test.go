package sign_test

import (
	"crypto/rsa"
	"testing"

	"github.com/arjanvaneersel/kit/sign"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var docHash sign.DocumentHash
var signature sign.Signature
var doc = struct {
	Name string
	Age  int
}{"Go Pher", 40}

func TestCreatePrivateKey(t *testing.T) {
	var err error
	privateKey, publicKey, err = sign.CreateKeyPair()
	if err != nil {
		t.Fatalf("[FAIL] Expected CreatePrivateKey to pass, but got error %q.", err)
	}
	t.Logf("[PASS] Private key: %v\nPublic key: %v.", privateKey, publicKey)
}

func TestHashDocument(t *testing.T) {
	var err error
	docHash, err = sign.HashDocument(&doc)
	if err != nil {
		t.Fatalf("[FAIL] Expected DocHash to pass, but got error %q.", err)
	}
	t.Logf("[PASS] Document hash: %s.", docHash)
}

func TestSign(t *testing.T) {
	var err error

	signature, err = sign.Sign(docHash, privateKey)
	if err != nil {
		t.Fatalf("[FAIL] Expected Sign to pass, but got error %q.", err)
	}
	t.Logf("[PASS] Signature: %s.", signature)
}

func TestValidSignature(t *testing.T) {
	if err := sign.ValidSignature(docHash, signature, publicKey); err != nil {
		t.Fatalf("[FAIL] Expected signature to be valid, but got error %q instead", err)
	}
}
