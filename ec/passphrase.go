package ec

import (
	"io"

	"github.com/SHDMT/btcec"
	"github.com/sammy00/bip38/encoding"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/text/unicode/norm"
)

// EncryptPassphrase generates a intermediate passphrase from the given one.
func EncryptPassphrase(rand io.Reader, passphrase string) (
	string, error) {

	var ownerEntropy [8]byte
	if _, err := rand.Read(ownerEntropy[:]); nil != err {
		return "", err
	}

	// ownerSalt=ownerEntropy, and pre->pass conversion if omitted
	pass, err := scrypt.Key(norm.NFC.Bytes([]byte(passphrase)), ownerEntropy[:],
		n1, r1, p1, keyLen1)
	if nil != err {
		return "", err
	}

	_, pub := btcec.PrivKeyFromBytes(btcec.S256(), pass)
	passPoint := pub.SerializeCompressed()

	return encoding.CheckEncode(
		noLotSequence[:], append(ownerEntropy[:], passPoint...)), nil
}

func EncryptPassphraseEx(rand io.Reader, passphrase string,
	lot, sequence uint32) (string, error) {
	var ownerEntropy [8]byte
	if _, err := rand.Read(ownerEntropy[:4]); nil != err {
		return "", err
	}

	ownerEntropy[4] = byte(lot >> 12 & 0xff)
	ownerEntropy[5] = byte(lot >> 4 & 0xff)
	ownerEntropy[6] = byte((lot << 4 & 0xf0) | (sequence >> 8 & 0x0f))
	ownerEntropy[7] = byte(sequence & 0xff)

	return "", nil
}