package crypt

import (
	"bytes"
	"crypto/md5"
	"io"
	"strings"

	"github.com/gtank/cryptopasta"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"rot13": func(s string) string { return strings.Map(rot13, s) },
	"md5":   md5Sum,
	"rndKey": func() []byte {
		k := cryptopasta.NewEncryptionKey()
		key := make([]byte, 32)
		copy(key, k[:])
		return key
	},
	"aes": func(key []byte, plain []byte) ([]byte, error) {
		var k [32]byte
		copy(k[:], key)
		return cryptopasta.Encrypt(plain, &k)
	},
	"unaes": func(key []byte, cipher []byte) ([]byte, error) {
		var k [32]byte
		copy(k[:], key)
		return cryptopasta.Decrypt(cipher, &k)
	},
}

func rot13(r rune) rune {
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
