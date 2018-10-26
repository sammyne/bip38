package bip38

import (
	"bytes"
	"crypto/aes"
	"errors"

	"golang.org/x/text/unicode/norm"

	"golang.org/x/crypto/scrypt"
)

func Decrypt(encrypted string, passphrase string) ([]byte, error) {
	// TODO: distinguish decoding routine based on version
	payload, _, err := CheckDecode(encrypted)
	if nil != err {
		return nil, err
	}

	dk, err := scrypt.Key(norm.NFC.Bytes([]byte(passphrase)),
		payload[:4], n, r, p, keyLen)
	if nil != err {
		return nil, err
	}

	C, err := aes.NewCipher(dk[32:])
	if nil != err {
		return nil, err
	}

	var plain [32]byte
	C.Decrypt(plain[:16], payload[4:20])
	C.Decrypt(plain[16:], payload[20:])

	priv := xor(plain[:], dk[:32])

	if !bytes.Equal(payload[:4], AddressHash(priv, false)) {
		return nil, errors.New("invalid address hash")
	}

	return priv, nil
}

// Encrypt encrypts the given private key byte sequence
// with the given passphrase
func Encrypt(data []byte, passphrase string) (string, error) {
	addrHash := AddressHash(data, false)

	dk, err := scrypt.Key(norm.NFC.Bytes([]byte(passphrase)),
		addrHash, n, r, p, keyLen)
	if nil != err {
		return "", err
	}

	var payload [36]byte
	copy(payload[:], addrHash) // append salt

	C, err := aes.NewCipher(dk[32:])
	if nil != err {
		return "", err
	}

	block := xor(data, dk[:32])
	C.Encrypt(payload[4:], block[:16])
	C.Encrypt(payload[20:], block[16:])

	return CheckEncode(payload[:], [3]byte{0x01, 0x42, 0xc0}), nil
}

// xor calculates the (x[0]^y[0], x[1]^y[1],..., x[32]^y[32])
func xor(x, y []byte) []byte {
	var out [32]byte
	for i := range out {
		out[i] = x[i] ^ y[i]
	}

	return out[:]
}