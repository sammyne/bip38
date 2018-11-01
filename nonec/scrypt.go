package nonec

import "golang.org/x/crypto/scrypt"

func init() {
	// dynamic checking of the proposed configuration for SCRYPT
	// so as to remove the redundant error checking later
	if _, err := scrypt.Key([]byte("bip38"), []byte("hello world"),
		N, R, P, KeyLen); nil != err {
		panic(err)
	}
}
