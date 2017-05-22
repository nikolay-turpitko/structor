package crypt

import (
	"bytes"
	"crypto/md5"
	"io"
	"strings"

	"github.com/gtank/cryptopasta"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// func rot13(s string) string
	// Performs simple rot13 obfuscation.
	"rot13": rot13,
	// func md5(r io.Reader) ([]byte, error)
	// Calculates md5 checksum of io.Reader's content and returns it as []byte.
	"md5": md5Sum,
	// func rndKey() []byte
	// Returns random encryption key for "aes"/"unaes".
	// See "github.com/gtank/cryptopasta".NewEncryptionKey().
	"rndKey": rndKey,
	// func aes(key, plain []byte) ([]byte, error)
	// Encrypts data using 256-bit AES-GCM. See "github.com/gtank/cryptopasta".Encrypt().
	"aes": aes,
	// func unaes(key, cipher []byte) ([]byte, error)
	// Decrypts data using 256-bit AES-GCM. See "github.com/gtank/cryptopasta".Decrypt().
	"unaes": unaes,
}

func rot13(s string) string { return strings.Map(mapRot13, s) }

func mapRot13(r rune) rune {
	if r >= 'a' && r <= 'z' {
		if r >= 'm' {
			return r - 13
		}
		return r + 13
	} else if r >= 'A' && r <= 'Z' {
		if r >= 'M' {
			return r - 13
		}
		return r + 13
	}
	return r
}

func md5Sum(r io.Reader) ([]byte, error) {
	h := md5.New()
	var b bytes.Buffer
	_, err := io.Copy(&b, r)
	if err != nil {
		return nil, err
	}
	_, err = h.Write(b.Bytes())
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func rndKey() []byte {
	k := cryptopasta.NewEncryptionKey()
	key := make([]byte, 32)
	copy(key, k[:])
	return key
}

func aes(key, plain []byte) ([]byte, error) {
	var k [32]byte
	copy(k[:], key)
	return cryptopasta.Encrypt(plain, &k)
}

func unaes(key, cipher []byte) ([]byte, error) {
	var k [32]byte
	copy(k[:], key)
	return cryptopasta.Decrypt(cipher, &k)
}
