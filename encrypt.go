package capitalcom

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func EncryptPassword(password string, key EncryptionKey) (string, error) {
	input := []byte(fmt.Sprintf("%s|%d", password, key.TimeStamp.UnixMilli()))
	base64Input := base64.StdEncoding.EncodeToString(input)

	keyBytes, err := base64.StdEncoding.DecodeString(key.EncryptionKey)
	if err != nil {
		return "", NewPasswordEncodingError(err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return "", NewPasswordEncodingError(err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", NewPasswordEncodingError(ErrPublicKeyTypeError)
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(base64Input))
	if err != nil {
		return "", NewPasswordEncodingError(err)
	}

	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)

	return encryptedBase64, nil
}
